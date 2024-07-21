package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	dbt "github.com/bmaayandexru/go_final_project/dbt"
)

const operationTypes string = "dywm"

type bmaErr struct {
	Error string `json:"error"`
}

type bmaId struct {
	Id string `json:"id"`
}

func NextDateHandle(res http.ResponseWriter, req *http.Request) {
	// лог-контроль
	fmt.Println("Получен запрос api/nextdate ")
	// запрос в строку
	s := fmt.Sprintf("Host: %s\nPath: %s\nMethod: %s", req.Host, req.URL.Path, req.Method)
	// лог-контроль
	fmt.Println(s)
	strNow := req.FormValue("now")
	strDate := req.FormValue("date")
	strRepeat := req.FormValue("repeat")
	fmt.Printf("now *%s* date *%s* repeat *%s*\n", strNow, strDate, strRepeat)
	now, err := time.Parse("20060102", strNow)
	if err != nil {
		return
	}
	retStr, err := NextDate(now, strDate, strRepeat)
	if err != nil {
		return
	}
	// отправка клиенту
	res.Write([]byte(retStr))
}

func bmaError(res http.ResponseWriter, sErr string, statuCode int) {
	var bE bmaErr
	bE.Error = sErr
	aBytes, _ := json.Marshal(bE)
	// *** лог контроль
	fmt.Println("bmaError: aBytes ", string(aBytes))
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statuCode)
	res.Write(aBytes)
}

func TaskHandle(res http.ResponseWriter, req *http.Request) {
	var task dbt.Task
	var buf bytes.Buffer
	// запрос в строку
	s := fmt.Sprintf("Получен запрос H: %s Path: %s M: %s", req.Host, req.URL.Path, req.Method)
	fmt.Println(s)

	switch req.Method {
	case "POST":
		// читаем тело запроса
		_, err := buf.ReadFrom(req.Body)
		if err != nil {
			// http.Error(res, err.Error(), http.StatusBadRequest)
			bmaError(res, fmt.Sprintf("TaskHandle: Ошибка чтения тела запроса: %s\n", err.Error()), http.StatusOK)
			return
		}
		// *** лог контроль
		fmt.Println("buf.Bytes():", buf.String())
		// десериализуем JSON в task
		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			//http.Error(res, err.Error(), http.StatusBadRequest)
			bmaError(res, fmt.Sprintf("TaskHandle: Ошибка десериализации: %s\n", err.Error()), http.StatusOK)
			return
		}
		// лог контроль
		fmt.Printf("Unmarshal Task Date *%s* Title *%s* Comment *%s* Repeat *%s*\n", task.Date, task.Title, task.Comment, task.Repeat)
		// анализ входных данных
		// Поле Title обязательное
		if len(task.Title) == 0 {
			bmaError(res, "TaskHandle: Поле `Задача*` пустое", http.StatusOK)
			return
		}
		// Если поле date не указано или содержит пустую строку,
		if len(task.Date) == 0 {
			// берётся сегодняшнее число.
			task.Date = time.Now().Format("20060102")
		} else {
			//  task.Date не пустое. пробуем распарсить
			_, err = time.Parse("20060102", task.Date)
			if err != nil {
				// ошибка разбора даты
				bmaError(res, fmt.Sprintf("TaskHandle: Ошибка разбора даты: %v\n", err), http.StatusOK)
				return
			}
		}
		// тут валидная строка в task.Date
		// это либо строка из текущей даты либо корректная строка
		// вызываем NextDate(time.Now(),task.Date,task.Repeat)
		// Если дата меньше сегодняшнего числа,  есть два варианта:
		if task.Date < time.Now().Format("20060102") {
			if len(task.Repeat) == 0 {
				// 1. если правило повторения не указано или равно пустой строке,
				// подставляется сегодняшнее число;
				task.Date = time.Now().Format("20060102")
			} else {
				// 2. при указанном правиле повторения вам нужно вычислить
				// и записать в таблицу дату выполнения, которая будет больше
				// сегодняшнего числа. Для этого используйте функцию NextDate(),
				// которую вы уже написали раньше.
				task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					bmaError(res, fmt.Sprintf("NextDate: %v", err), http.StatusOK)
					return
				}
			}
		} else {
			task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
			if err != nil {
				bmaError(res, fmt.Sprintf("NextDate: %v", err), http.StatusOK)
				return
			}
		}
		fmt.Println("Добавлена в базу ", task)
		// Task положить в базу и определить id
		resDBEx, err := dbt.AddTask(task)
		if err != nil {
			// оборачиваем ошибку в свою
			bmaError(res, fmt.Sprintf("TaskHandle: Ошибка при добавлении в БД: %v\n", err), http.StatusOK)
			return
		}
		// определить id
		id, err := resDBEx.LastInsertId()
		if err != nil {
			bmaError(res, fmt.Sprintf("TaskHandle: Ошибка LastInsetId(): %v\n", err), http.StatusOK)
			return
		}
		var bId bmaId
		bId.Id = strconv.Itoa(int(id))
		// вернуть id из базы в json
		// кодируем id для отправки ответа
		arrBytes, err := json.Marshal(bId)
		if err != nil {
			bmaError(res, fmt.Sprintf("TaskHandle: Ошибка json.Marshal(id): %v\n", err), http.StatusOK)
			return
		}
		// *** лог контроль
		fmt.Printf("ret json *%s*\n", string(arrBytes))
		// запись результата в JSON
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusOK)
		res.Write(arrBytes)

	case "GET":
	case "DELETE":
	}
}

// фнкция пересчета следующей даты
// now — время от которого ищется ближайшая дата;
// date — строка времени в формате 20060102, от которого начинается отсчёт повторений;
// repeat — правило повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {
	template := "20060102"
	// fmt.Printf("ND: now %s, date *%s*, repeat *%s*\n", now.Format(template), date, repeat)
	startDate, err := time.Parse(template, date)
	if err != nil {
		// ошибка в стартовой дате
		return "", err
	}
	// разложить repeat в слайс строк
	repSlice := strings.Split(repeat, " ")
	// fmt.Println("ND: slice", repSlice)
	if len(repSlice[0]) == 0 {
		// поле repeat не задано
		if now.Before(startDate) {
			return date, nil
		} else {
			return now.Format(template), nil
		}
	}
	// тут слайс не пустой. проверяем певый элемент на соответствие
	// первая строка - тип повторения. один символ из стоки "dywm"
	if len(repSlice[0]) != 1 {
		return "", errors.New("Длина типа не равна 1")
	}
	if !strings.Contains(operationTypes, repSlice[0]) {
		return "", errors.New("Неизвестный тип операции повторения")
	}
	switch repSlice[0] {
	case "d": // d дни
		// fmt.Println("дни")
		// у d может быть только одно число
		// проверить есть ли параметры за "d"
		if len(repSlice) < 2 {
			return "", errors.New("d: нет указания дней")
		}
		if len(repSlice) > 2 {
			return "", errors.New("d: много параметров")
		}
		// разложить rs[1] на слайс
		repSlice1 := strings.Split(repSlice[1], ",")
		if len(repSlice1) != 1 {
			return "", errors.New("d: число дней указано не одним числом")
		}
		dcount, err := strconv.Atoi(repSlice1[0])
		if err != nil {
			return "", errors.New("d: параметр не число") // параметр не число
		}
		// число от 1 до 400 включительно
		if (dcount < 1) || (dcount > 400) {
			return "", errors.New("d: число вне диапазона (<1 >400)") // число вне диапазона
		}
		// тут всё корректно. можно возвращать значение
		/*
			if dcount == 1 && startDate.Before(now) {
				// не проходит тест поэтому так
				return now.Format(template), nil
				// тест прошел. очень странно
			}
		*/
		/*
			itc := 1
			if startDate.Before(now) {
				// расчет полных периодов до нужной даты
				itc = int(now.Sub(startDate).Hours()/24)/dcount + 1
			}
			startDate = startDate.AddDate(0, 0, dcount*itc)
		*/
		/*
			for startDate.Before(now) {
				startDate = startDate.AddDate(0, 0, dcount)
			}
		*/
		for startDate.Format(template) < now.Format(template) {
			startDate = startDate.AddDate(0, 0, dcount)
		}
		return startDate.Format(template), nil

	case "y": // y год
		// !!! в любом случае идет перенос даты на год хотя бы однократно.
		// дальше сравниваем получившуюся даты с текущей
		// и если получившаяся меньше добавляем еще год
		// до тех пор пока получишаяся не будет больше текущей
		// у Y не может быть 2го слайса/числа
		if len(repSlice) != 1 {
			return "", errors.New("y: количество параметров != 0") // ошибка количества параметров
		}
		startDate = startDate.AddDate(1, 0, 0)
		for startDate.Before(now) {
			startDate = startDate.AddDate(1, 0, 0)
		}
		return startDate.Format(template), nil

	case "w":
		// w дни недели
		// у w может быть слайс из 1-7
		if len(repSlice) != 2 {
			return "", errors.New("w: количество параметров != 2") // ошибка количества параметров
		}
		repSlice1 := strings.Split(repSlice[1], ",") // второй параметр в слайс строк
		for i := 0; i < len(repSlice1); i++ {        // проверка на допустимые значения
			if (repSlice1[i] < "1") || (repSlice1[i] > "7") { // если вне диапазона
				return "", errors.New("w: один из параметров за пределами диапазона (<1 >7)") // ошибка
			}
		}
		// мапа из значений с преобразование в int'ы а затем в Weedday
		mapWeekDays := make(map[time.Weekday]bool)
		for _, strDay := range repSlice1 {
			iDay, _ := strconv.Atoi(strDay)
			if iDay == 7 {
				iDay = 0 // вс теперь 0 а не 7
			}
			mapWeekDays[time.Weekday(iDay)] = true
		}
		curDay := now
		curDay = curDay.AddDate(0, 0, 1)
		_, found := mapWeekDays[curDay.Weekday()]
		for !found {
			curDay = curDay.AddDate(0, 0, 1)
			fmt.Println(curDay)
			_, found = mapWeekDays[curDay.Weekday()]
		}
		return curDay.Format(template), nil
	case "m":
		// m дни месяца
		sDate, err := NextDateMonth(now, startDate, repSlice)
		return sDate, err
	} // switch
	return "", errors.New("не удалось определить следующую дату")
}

func makeSlice3Month(date time.Time) []time.Time {
	retSl := make([]time.Time, 0, 3)
	retSl = append(retSl, time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local))
	date = date.AddDate(0, 1, 0)
	retSl = append(retSl, time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local))
	date = date.AddDate(0, 1, 0)
	retSl = append(retSl, time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local))
	return retSl
}

func checkDay(year int, month int, day int) bool {
	// месяцы в которых нет 31го числа
	msless31 := map[int]bool{2: true, 4: true, 6: true, 9: true, 11: true}

	if day >= 1 && day <= 28 {
		// дни [1..28] проходят проверку
		return true
	} else {
		// проверка наличия в месяце 29, 30, 31 чисел
		if day == 31 && msless31[month] {
			// 31е число есть не во всех месяцах
			return false
		}
		if year%4 == 0 {
			// високосный год
			if month == 2 && day == 30 {
				return false
			}
		} else {
			// не високосный год
			if month == 2 && (day == 29 || day == 30) {
				return false
			}
		}
		// проверка прошла возвращаем true
		return true
	}
}

func NextDateMonth(now time.Time, startDate time.Time, repSlice []string) (string, error) {
	// template := "20060102"

	if len(repSlice) < 2 {
		return "", errors.New("m: не узазана дата/даты месяца")
	}
	// проверка списка первого параметра [1..31,-1,-2]
	repSlice1 := strings.Split(repSlice[1], ",") // из 2й группы делаем слас строк
	sliDays := make([]int, 0, len(repSlice1))    // создаем слайс дней
	for _, strDay := range repSlice1 {           // []string --> []int
		iDay, err := strconv.Atoi(strDay)
		if err == nil {
			if ((iDay >= 1) && (iDay <= 31)) || iDay == -1 || iDay == -2 { // если в диапазане
				sliDays = append(sliDays, iDay)
			} else {
				return "", errors.New("m: день месяца вне диапазона") // ошибка
			}
		} else {
			return "", errors.New(fmt.Sprintf("m: День указан не числом. Ошибка:%v \n", err))
		}
	}
	sort.Ints(sliDays)        // в slDays отсортированные дни
	var sldMonths []time.Time // слайс первых чисел нужных месяцев
	if len(repSlice) == 2 {
		// тут только один список. создаем слайс из трёх месяцев
		if startDate.Before(now) {
			sldMonths = makeSlice3Month(now)
		} else {
			sldMonths = makeSlice3Month(startDate)
		}
	} else {
		// тут два списка len >= 3
		// должны получить слайс дат для месяцев из второго списка
		// проверяем второй список на корректность [1..12]
		repSlice2 := strings.Split(repSlice[2], ",") // третий параметр в слайс строк
		sliMonths := make([]int, 0, len(repSlice2))  // слайс целых чисел месяцев
		for _, strMonth := range repSlice2 {
			iMonth, err := strconv.Atoi(strMonth)
			if err == nil {
				if (iMonth >= 1) && (iMonth <= 12) { // если в диапазоне
					sliMonths = append(sliMonths, iMonth) // добавляем
				} else {
					return "", errors.New("m: месяц за пределами диапазона") // ошибка
				}
			} else {
				return "", errors.New("m: указано не число") // в слайсе не число
			}
		}
		sort.Ints(sliMonths) // сортировка
		// создаем слайс дат для текущего и следующего года
		sldYears := make([]time.Time, 0)
		sldYears = append(sldYears, now)
		sldYears = append(sldYears, now.AddDate(1, 0, 0))
		// создаем слайс из дат для каждого года
		for _, dYear := range sldYears {
			for _, iMonth := range sliMonths {
				// добавляем в slMonths дату (dYear,iMonth,01)
				sldMonths = append(sldMonths, time.Date(dYear.Year(), time.Month(iMonth), 1, 0, 0, 0, 0, time.Local))
			}
		}
	}
	// тут сфоормирован слайс из дат для месяцев и годов
	// формируем слайс строк дат по дням, т к требуется сортировка
	slsDays := make([]string, 0)
	for _, dMonth := range sldMonths {
		for _, iDay := range sliDays {
			// создаем список конкретных
			// из даты dMonth берем месяц и год
			// проверяем есть ли iDay в этом месяце
			if checkDay(dMonth.Year(), int(dMonth.Month()), iDay) {
				// если есть формируем дату dMonth.Year(), dMonth.Month(). iDay
				// добавляем строкой в список дней, который потом отсортируем и выберем нужный
				slsDays = append(slsDays, time.Date(dMonth.Year(), dMonth.Month(), iDay, 0, 0, 0, 0, time.Local).Format("20060102"))
			}
		}
	}
	// сортируем список из строк дат
	if len(slsDays) == 0 {
		return "", errors.New("m: Не возможно определить дату по указанным параметрам (m 30,31 2)")
	}
	sort.Strings(slsDays)
	// выбираем нужный день и возвращаем
	for _, sDay := range slsDays {
		dDay, _ := time.Parse("20060102", sDay)
		if now.Before(dDay) {
			return sDay, nil
		}
	}
	// если список пуст то это m 30,31 2
	return "", errors.New("m: Не возможно определить дату по указанным параметрам (m 30,31 2)")
}
