package handler

type Translator interface {
	Translate(text, sourceLang, targetLang string) (string, error)
}

type Handler struct {
	Tran Translator
}

func NewHandler(t Translator) *Handler {
	return &Handler{
		Tran: t,
	}
}

func (h *Handler) Process(text, sourceLang, targetLang string) (string, error) {
	result, err := h.Tran.Translate(text, sourceLang, targetLang)
	if err != nil {
		return "", err
	}

	return result, nil
}
