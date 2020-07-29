package product

type Product struct {
	Id           int     `json:"id"`
	ProductCode  string  `json:"productCode"`
	ProductName  string  `json:"productName"`
	Description  string  `json:"description"`
	StandardCost float64 `json:"standardCost"`
	ListPrice    float64 `json:"listPrice"`
	Category     string  `json:"category"`
}

type ProductList struct {
	Data         []*Product `json:"data"`
	TotalRecords int        `json:"TotalRecords"`
}

type ProductTop struct {
	Id          int     `json:"id"`
	ProductName string  `json:"productName"`
	Sold        float64 `json:"sold"`
}

type ProductTopResponse struct {
	Data  []*ProductTop `json:"data"`
	Total float64       `json:"Total"`
}
