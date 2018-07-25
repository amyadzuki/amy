package widget

type Performance struct {
	Panel *gui.Panel
	Value, Units *gui.Label
}

func (w *Performance) NewPerformance(w [2]*gui.Label, large int, label string) {
	w.Panel = gui.NewPanel(0, 0)
	w.Panel.SetLayout(gui.NewDockLayout())
	w.Units = gui.NewLabel(label)
	w.Units.SetLayoutParams(&gui.DockLayoutParams{gui.DockRight})
	w.Value = gui.NewLabel(strconv.Itoa(large))
	w.Value.SetLayoutParams(&gui.DockLayoutParams{gui.DockRight})
	w.Panel.Add(w.Units)
	w.Panel.Add(w.Value)
	w.Units.SetText("")
}
