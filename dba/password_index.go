package dba

import (
	"log"
	"strconv"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/pborman/uuid"
)

// PasswordIndex is a minimal set of entry points into bleve.Index
// that we need to operate. Brought out to simplify mocking.
// See bleve.Index for documentation and usage.
type PasswordIndex interface {
	// Inserts the provided data bundle into the password database.
	Index(id string, data interface{}) error

	// Searches the password database for the specified search string
	// and returns a search result.
	Search(req *bleve.SearchRequest) (*bleve.SearchResult, error)

	// ListUsers searches the password database for the specified match string and
	// returns a list of users of not more than size entries starting from
	// from in the result lis or error.
	// TODO(rjk): Figure out if I want to specify the fields.
	ListUsers(userquery string, size, from int) ([]map[string]interface{}, error)

	// Returns the document map for a given document (i.e. password entry)
	// when presented with the entry's UUID
	// TODO(rjk): Make sure that the returned map can be indexed.
	// And maybe give these better names. And use UUIDs for type clarity.
	MapForDocument(id string) (map[string]interface{}, error)

	// Deletes the specificed password entry.
	Delete(id string) error
}

type blevePasswordIndex struct {
	idx bleve.Index
}

func (pdoc *blevePasswordIndex) MapForDocument(id string) (map[string]interface{}, error) {
	idx := pdoc.idx
	doc, err := idx.Document(id)
	if err != nil {
		return nil, err
	}
	return MakeMapFromDocument(doc)
}

// TODO(rjk): Fix the layering violation of leaking bleve into the interface.
// Update: ListUsers resolves this in part.
func (pdoc *blevePasswordIndex) Search(req *bleve.SearchRequest) (*bleve.SearchResult, error) {
	return pdoc.idx.Search(req)
}

func (pdoc *blevePasswordIndex) ListUsers(userquery string, size, from int) ([]map[string]interface{}, error) {
	log.Println("called ListUsers", userquery, size, from)

	var queryop query.Query
	if userquery == "" {
		queryop = bleve.NewMatchAllQuery()
	} else {
		queryop = bleve.NewWildcardQuery(userquery)
	}

	sreq := bleve.NewSearchRequest(queryop)
	sreq.Fields = []string{"name", "role", "display_name"}
	sreq.Size = size
	sreq.From = from

	rawresults, err := pdoc.idx.Search(sreq)
	if err != nil {
		log.Println("error in search", err)
		return nil, err
	}

	if len(rawresults.Hits) < 1 {
		log.Println("no users in search")
		return make([]map[string]interface{}, 0, 0), nil
	}
	users := make([]map[string]interface{}, 0, len(rawresults.Hits))

	for i, sr := range rawresults.Hits {
		u := make(map[string]interface{})
		for k, v := range sr.Fields {
			// Could test and drop the unfortunate?
			u[k] = v.(string)
		}

		uuidcasted := uuid.UUID(sr.ID)
		// I thought about encrypting the UUIDs. But to get this content, one
		// must already have the admin role and that is enforced server side
		// via a strongly encrypted cookie. And they are cryptographically
		// difficult to guess already.
		u["uuid"] = uuidcasted.String()

		// It's conceivable that offsetting this by where we are in the list is foolish.
		u["index"] = strconv.FormatInt(int64(i+from), 10)
		users = append(users, u)
	}
	return users, nil
}

func (pdoc *blevePasswordIndex) Delete(id string) error {
	return pdoc.idx.Delete(id)
}

func (pdoc *blevePasswordIndex) Index(id string, data interface{}) error {
	return pdoc.idx.Index(id, data)
}

func OpenPassword(persistentroot string) (PasswordIndex, error) {
	passwordfile, err := OpenBleve(persistentroot, PasswordFile)
	if err != nil {
		return nil, err
	}
	return &blevePasswordIndex{passwordfile}, err
}
