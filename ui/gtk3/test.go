package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
)

func main() {
	// Инициализируем GTK.
	gtk.Init(nil)

	// Создаём билдер
	b, err := gtk.BuilderNew()
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Загружаем в билдер окно из файла Glade
	err = b.AddFromFile("/Users/mac/工具.glade")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Получаем объект главного окна по ID
	obj, err := b.GetObject("window_main")
	if err != nil {
		log.Fatal("Ошибка:", err)
	}

	// Преобразуем из объекта именно окно типа gtk.Window
	// и соединяем с сигналом "destroy" чтобы можно было закрыть
	// приложение при закрытии окна
	win := obj.(*gtk.Window)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	object,err:=b.GetObject("bbQueryOrder")
	if err!=nil{
		panic(err.Error())
	}
	query:=object.(*gtk.MenuItem)
	query.Connect("activate",func(){
		object,err=b.GetObject("orderId")
		dialog:=object.(*gtk.Dialog)
		dialog.SetTransientFor(win)
		dialog.Show()
	})
	// Отображаем все виджеты в окне
	win.ShowAll()

	// Выполняем главный цикл GTK (для отрисовки). Он остановится когда
	// выполнится gtk.MainQuit()
	gtk.Main()
}