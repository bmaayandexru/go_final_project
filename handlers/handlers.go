package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const operationTypes string = "dywm"

// отладочный обработчик. убрать
func DbgHandle(res http.ResponseWriter, req *http.Request) {
	// лог-контроль
	fmt.Println("Получен запрос /d ")
	// запрос в строку
	s := fmt.Sprintf("Host: %s\nPath: %s\nMethod: %s", req.Host, req.URL.Path, req.Method)
	// лог-контроль
	fmt.Println(s)
	// отправка клиенту
	res.Write([]byte(s))
}

func NextDateHandle(res http.ResponseWriter, req *http.Request) {
	// лог-контроль
	fmt.Println("Получен запрос api/nextdate ")
	// запрос в строку
	s := fmt.Sprintf("Host: %s\nPath: %s\nMethod: %s", req.Host, req.URL.Path, req.Method)
	// лог-контроль
	fmt.Println(s)
	//rs := req.FormValue()
	//fmt.Println(rs)
	// отправка клиенту
	res.Write([]byte(s))
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
	if len(repSlice) == 0 {
		// слайс rs пустой
		return "", errors.New("repeat is empty")
	}
	// тут слайс не пустой. проверяем певый элемент на соответствие
	// первая строка - тип пересчета. один символ из стоки "dywm"
	if len(repSlice[0]) != 1 {
		// неправильная длина типа
		return "", errors.New("repetition lenght error")
	}

	if !strings.Contains(operationTypes, repSlice[0]) {
		// неизвестный тип операции повторения
		return "", errors.New("repetition type error")
	}
	switch repSlice[0] {
	case "d": // d дни
		// fmt.Println("дни")
		// у d может быть только одно число
		// проверить есть ли параметры за "d"
		if len(repSlice) < 2 {
			return "", errors.New("d: There is no number of days") // нет указания дней
		}
		if len(repSlice) > 2 {
			return "", errors.New("d: many params") // много параметров
		}
		// разложить rs[1] на слайс
		repSlice1 := strings.Split(repSlice[1], ",")
		if len(repSlice1) != 1 {
			return "", errors.New("d: param count error (!=1)") // число не одно
		}
		dcount, err := strconv.Atoi(repSlice1[0])
		if err != nil {
			return "", errors.New("d: param is not number") // параметр не число
		}
		// число от 1 до 400 включительно
		if (dcount < 1) || (dcount > 400) {
			return "", errors.New("d: value error (<1 >400)") // число вне диапазона
		}
		// тут всё корректно. можно возвращать значение
		// fmt.Printf("кол-во дней %d \n", dcount)                     // *** лог контроль
		// fmt.Printf("Дата старта %s \n", startDate.Format(template)) // *** лог контроль
		// расчет полных периодов до нужной даты
		itc := int(now.Sub(startDate).Hours()/24)/dcount + 1
		// fmt.Printf("Полных периодов +1 %v \n", itc)
		startDate = startDate.AddDate(0, 0, dcount*itc)
		// fmt.Printf("Дата на выдачу %s\n", startDate.Format(template)) // *** лог контроль
		return startDate.Format(template), nil

	case "y": // y год
		// fmt.Println("год")
		// у Y не может быть 2го слайса/числа
		// разложить rs[1] на слайс
		if len(repSlice) != 1 {
			return "", errors.New("y: param count error (!= 0)") // ошибка количества параметров
		}
		// fmt.Printf("Дата старта %s\n", startDate.Format(template))
		for startDate.Before(now) {
			startDate = startDate.AddDate(1, 0, 0)
			fmt.Printf("Итерация %v\n ", startDate)
		}
		// fmt.Printf("Дата на выдачу %s\n", startDate.Format(template))
		return startDate.Format(template), nil

	case "w":
		// w дни недели
		// у w может быть слайс из 1-7
		// fmt.Println("неделя")
		if len(repSlice) != 2 {
			return "", errors.New("w: paramcount error (!=2)") // ошибка количества параметров
		}
		repSlice1 := strings.Split(repSlice[1], ",") // второй параметр в слайс строк
		// fmt.Println(repSlice1)                       // *** лог контроль
		for i := 0; i < len(repSlice1); i++ { // проверка на допустимые значения
			if (repSlice1[i] < "1") || (repSlice1[i] > "7") { // если вне диапазона
				return "", errors.New("w: day out of range (<1 >7)") // ошибка
			}
		}
		mapWeekDays := make(map[time.Weekday]bool) // мапа из значений с преобразование в int'ы а затем в Weedday
		for _, strDay := range repSlice1 {
			iDay, _ := strconv.Atoi(strDay)
			if iDay == 7 {
				iDay = 0 // вс теперь 0 а не 7
			}
			mapWeekDays[time.Weekday(iDay)] = true
		}
		// fmt.Println(mapWeekDays) // *** лог контроль
		curDay := now
		curDay = curDay.AddDate(0, 0, 1)
		// fmt.Println(curDay)
		_, found := mapWeekDays[curDay.Weekday()]
		for !found {
			curDay = curDay.AddDate(0, 0, 1)
			fmt.Println(curDay)
			_, found = mapWeekDays[curDay.Weekday()]
		}
		return curDay.Format(template), nil
	case "m":
		// m дни месяца
		// у m может быть слайс из чисел 1..31,-1,-2
		// fmt.Println("месяц")
		if len(repSlice) < 2 {
			return "", errors.New("m: There is no number of days")
		}
		// проверка списка первого параметра [1..31,-1,-2]
		repSlice1 := strings.Split(repSlice[1], ",") // из 2й группы делаем слас строк
		// fmt.Println(repSlice1)                       // *** лог контроль
		slDays := make([]int, 0, len(repSlice1)) // создаем слайс дней
		for _, strDay := range repSlice1 {       // []string --> []int
			iDay, err := strconv.Atoi(strDay)
			if err == nil {
				if ((iDay >= 1) && (iDay <= 31)) || iDay == -1 || iDay == -2 { // если в диапазане
					slDays = append(slDays, iDay)
				} else { // иначе
					return "", errors.New("m: A day out of range") // ошибка
				}
			}
		}
		sort.Ints(slDays) // в slDays отсортированные дни
		// fmt.Println("Дни ", slDays) // *** лог контроль
		if len(repSlice) == 2 {
			// тут только один список
			// формируем даты на текущий и следующий месяц
			nowDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local) // 1е число текущего месяца
			nextDate := nowDate.AddDate(0, 1, 0)                                     // 1e число следующего месяца
			// проверяем значения слайса дней на присутствие в каждом меняце nowDate и nextDate
			// от now и next Берутся только год и месяц
			if err := checkDaysMonth(slDays, nowDate); err != nil {
				return "", err
			}
			if err := checkDaysMonth(slDays, nextDate); err != nil {
				return "", err
			}
			// строки дат в формате "20060102"
			slDates := make([]string, 0, len(slDays)*2)
			for _, day := range slDays {
				if day > 0 {
					// добавляем с слайс даты для каждого месяца
					slDates = append(slDates, time.Date(nowDate.Year(), nowDate.Month(), day, 0, 0, 0, 0, time.Local).Format(template))
					slDates = append(slDates, time.Date(nextDate.Year(), nextDate.Month(), day, 0, 0, 0, 0, time.Local).Format(template))
				} else {
					// day < 0
					// взять 1е число следующего месяца и сложить с day с помощья AddDate(0,0,day)
					nowD := time.Date(nowDate.Year(), nowDate.Month(), 1, 0, 0, 0, 0, time.Local)
					nowD = nowD.AddDate(0, 1, day)
					slDates = append(slDates, nowD.Format(template))
					nextD := time.Date(nextDate.Year(), nextDate.Month(), 1, 0, 0, 0, 0, time.Local)
					nextD = nextD.AddDate(0, 1, day)
					slDates = append(slDates, nextD.Format(template))
				}
			}
			sort.Strings(slDates)
			// fmt.Println("Даты")
			//for _, date := range slDates { // *** лог контроль
			//	fmt.Println(date)
			//}
			// ищем 1ю дату, которая больше now. это и будет результат
			for _, date := range slDates {
				t, e := time.Parse(template, date)
				if e == nil {
					fmt.Println("Сравнение ", date, now.Format(template))
					if now.Before(t) {
						return date, nil
					} else {
						fmt.Printf("Drop date %s\n", date)
					}
				}
			}
		} else {
			// тут два списка len >= 3
			// проверяем второй список на корректность [1..12]
			repSlice2 := strings.Split(repSlice[2], ",") // третий параметр в слайс строк
			//fmt.Println(repSlice2)                       // *** лог контроль
			slMonth := make([]int, 0, len(repSlice2)) // мапа месяцев
			for _, strMonth := range repSlice2 {      // слайс в мапу
				iMonth, err := strconv.Atoi(strMonth)
				if err == nil {
					if (iMonth >= 1) && (iMonth <= 12) { // если в диапазоне
						slMonth = append(slMonth, iMonth) // добавляем
					} else { // иначе
						return "", errors.New("m: month out of range") // ошибка
					}
				} else {
					return "", errors.New("m: The specified number is not") // в слайсе не число
				}
			}
			sort.Ints(slMonth) // сортировка
			// for i, m := range slMonth { // *** лог контроль
			//	fmt.Println(i, m)
			// }
			// определяем текущий и следующий годы
			nowDate := now
			nextDate := now.AddDate(1, 0, 0)
			// проверяем числа меяцев в каждом году
			for _, month := range slMonth {
				if err := checkDaysMonth(slDays, time.Date(nowDate.Year(), time.Month(month), 1, 0, 0, 0, 0, time.Local)); err != nil {
					return "", err
				}
				if err := checkDaysMonth(slDays, time.Date(nextDate.Year(), time.Month(month), 1, 0, 0, 0, 0, time.Local)); err != nil {
					return "", err
				}
			}
			// формируем слайс из дат за текущий и следующий год
			slDates := make([]string, 0, len(slMonth)*len(slDays)*2)
			for _, month := range slMonth {
				for _, day := range slDays {
					if day > 0 {
						slDates = append(slDates, time.Date(nowDate.Year(), time.Month(month), day, 0, 0, 0, 0, time.Local).Format(template))
						slDates = append(slDates, time.Date(nextDate.Year(), time.Month(month), day, 0, 0, 0, 0, time.Local).Format(template))
					} else {
						// day < 0
						nowD := time.Date(nowDate.Year(), time.Month(month), 1, 0, 0, 0, 0, time.Local)
						nowD = nowD.AddDate(0, 1, day)
						slDates = append(slDates, nowD.Format(template))
						nextD := time.Date(nextDate.Year(), time.Month(month), 1, 0, 0, 0, 0, time.Local)
						nextD = nextD.AddDate(0, 1, day)
						slDates = append(slDates, nextD.Format(template))
					}
				}
			}
			sort.Strings(slDates)
			// fmt.Println("Даты")
			// for _, date := range slDates { // *** лог контроль
			// 	fmt.Println(date)
			// }
			// ищем 1ю дату, которая больше now. это и будет результат
			for _, date := range slDates {
				t, e := time.Parse(template, date)
				if e == nil {
					//	fmt.Println("Сравнение ", date, now.Format(template))
					if now.Before(t) {
						return date, nil
					} else {
						fmt.Printf("Drop date %s\n", date)
					}
				}
			}
		}
	} // switch
	return "", errors.New("ND: end func")
}

// функция checkDaysMonth проверяет дни на наличие в указанном месяце month в указанном году year
// если ошибок нет, возвращается error == nil
func checkDaysMonth(slDays []int, date time.Time) error {
	// месяцы в которых нет 31го числа
	msless31 := map[int]bool{2: true, 4: true, 6: true, 9: true, 11: true}
	for _, day := range slDays {
		if day <= 28 {
			// все дни до 28 числа включительно проходят проверку
			continue
		}
		// не все 31е числа есть в месяцах
		if day == 31 && msless31[int(date.Month())] {
			return errors.New(fmt.Sprintf("cDM error: There is no %dth in the %d month of %d", day, date.Month(), date.Year()))
		}
		// февраль проверяем отдельно
		// если високосный проверяем 30
		// если не високосный проверяем 29 30
		if date.Month() == 2 {
			if date.Year()%4 == 0 {
				// високосный год
				if day == 30 {
					return errors.New(fmt.Sprintf("cDM error: There is no %dth in the %d month of %d", day, date.Month(), date.Year()))
				}
			} else {
				// не високосный год
				if day == 30 || day == 29 {
					return errors.New(fmt.Sprintf("cDM error: There is no %dth in the %d month of %d", day, date.Month(), date.Year()))
				}
			}
		}
	}
	return nil
}
