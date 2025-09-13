package flags

type Index string

const (
	Relevance Index = "relevance"
	Downloads Index = "downloads"
	Follows   Index = "follows"
	Newest    Index = "newest"
	Updated   Index = "updated"
)

var ValidIndex = []string{
	Relevance.String(),
	Downloads.String(),
	Follows.String(),
	Newest.String(),
	Updated.String(),
}

func (i Index) String() string {
	return string(i)
}
