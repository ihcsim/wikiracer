package wikiracer

type mockWiki struct {
	pages map[string]*Page
}

func NewMockWiki() *mockWiki {
	testData := map[string]*Page{
		"Alexander the Great": &Page{ID: 1000, Title: "Alexander the Great", Namespace: 0, Links: []string{"Greek language"}},
		"Fruit anatomy":       &Page{ID: 1001, Title: "Fruit anatomy", Namespace: 0, Links: []string{"Segment"}},
		"Greek language":      &Page{ID: 1002, Title: "Greek language", Namespace: 0, Links: []string{"Fruit anatomy"}},
		"Mike Tyson":          &Page{ID: 1003, Title: "Mike Tyson", Namespace: 0, Links: []string{"Alexander the Great"}},
		"Segment":             &Page{ID: 1004, Title: "Segment", Namespace: 0, Links: []string{}},
	}
	return &mockWiki{pages: testData}
}

func (m *mockWiki) FindPage(title string) (*Page, error) {
	page, exist := m.pages[title]
	if !exist {
		return nil, PageNotFound{page}
	}

	return page, nil
}
