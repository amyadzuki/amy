package widget

import (
	"github.com/g3n/engine/gui"
)

type Small struct {
	Panel *gui.Panel
	Label *gui.Label
}

func NewSmall(label string) (w *Small) {
	w = new(Small)
	w.Init(label)
	return
}

func (w *Small) Init(label string) {
	w.Panel = gui.NewPanel(0, 0)
	w.Label = gui.NewLabel(label)
	w.Panel.Add(w.Label)
	lw, lh := float64(w.Label.TotalWidth()), float64(w.Label.TotalHeight())
	w.Panel.SetWidth(float32(lw))
	w.Panel.SetHeight(float32(lh))
}
