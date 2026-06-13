package click

// Click позволяет получить ID трекерного клика.
type Click interface {
	ID() string     // ID клика в трекере.
	String() string // String возвращает строковое значение клика.
}

// Tracker возвращает информацио о трекере.
type Tracker interface {
	ID() int64        // ID возвращает локальный идентификатор трекера.
	Host() string     // Host возврщает http/https хост трекера.
	ClickURL() string // ClickURL возвращает URL путь для передачи клика в трекер.
	String() string   // String вовзращает строковое значение трекера.
}

// RoutedClick позволяет использовать интерфейс Click,
// а так же содержит данные о том из какого трекера пришел клик.
type RoutedClick interface {
	Click
	Tracker() Tracker
}
