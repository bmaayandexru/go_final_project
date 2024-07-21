package tests

var Port = 7540

// путь к базе для тестов
var DBFile = "../scheduler.db"

// путь к базе для запуска
// эту переменную использую у себя в коде
var DBFileRun = "scheduler.db"

// если использовать одну и ту же переменную и в тестах и в коде, то тесты не проходят
// не находит файл с базой

// var FullNextDate = false
var FullNextDate = true
var Search = false
var Token = ``
