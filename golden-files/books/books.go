package books

type Book struct {
	ISBN      string
	Title     string
	Author    string
	Pages     int
	Publisher string
	Price     int
}

var Books []Book = []Book{
	Book{
		ISBN:      "978-1591847786",
		Title:     "Hooked",
		Author:    "Nir Eyal",
		Pages:     256,
		Publisher: "Portfolio",
		Price:     19,
	},
	Book{
		ISBN:      "978-1434442017",
		Title:     "The Great Gatsby",
		Author:    "F. Scott Fitzgerald",
		Pages:     140,
		Publisher: "Wildside Press",
		Price:     12,
	},
	Book{
		ISBN:      "978-1784756260",
		Title:     "Then She Was Gone: A Novel",
		Author:    "Lisa Jewell",
		Pages:     448,
		Publisher: "Arrow",
		Price:     29,
	},
	Book{
		ISBN:      "978-1094400648",
		Title:     "Think Like a Billionaire",
		Author:    "James Altucher",
		Pages:     852,
		Publisher: "Scribd, Inc.",
		Price:     9,
	},
}
