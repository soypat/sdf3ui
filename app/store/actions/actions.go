package actions

// Context contains all information to completely define a page layout.
// It also has information on last page visited to enable the use of a "back" button.
// It is not an "action" item strictly speaking.
type Context struct {
	// Defines the current page.
	Page Page
	// Referrer contains information on last page visited.
	Referrer *Context
	// Action can contain the executed action struct along with all
	// data contained. Can be very useful though author is on the
	// fence on whether it is good practice.
	// Action interface{} // Uncomment for use.
}

type Page int

const (
	PageLanding Page = iota
	PageNewItem
)

// PageSelect Navigates view to new page.
type PageSelect struct {
	Page Page
}

// GetShape action.
type GetShape struct{}

// Back button pressed. Navigate to previous page.
type Back struct{}

// Refresh just updates page by calling rendering functions.
type Refresh struct{}
