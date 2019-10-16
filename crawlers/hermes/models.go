package hermes

const (
	searchBaseURL = "https://bck.hermes.cn/product?locale=cn_zh&sort=relevance&searchterm="
	detailBaseURL = "https://www.hermes.cn/cn/zh"
)

type ProductResponse struct {
	Total    int64      `json:"total"`
	Products []*Product `json:"products"`
}

// Product Product represent single product info
type Product struct {
	SKU   string `json:"sku"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (p *Product) GetDetailURL() string {
	return detailBaseURL + p.URL
}
