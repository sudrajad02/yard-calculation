package schemas

// Request
type SuggestContainerRequest struct {
	Yard            string  `json:"yard"`
	ContainerNumber string  `json:"container_number"`
	ContainerSize   int     `json:"container_size"`
	ContainerHeight float64 `json:"container_height"`
	ContainerType   string  `json:"container_type"`
}

type PlaceContainerRequest struct {
	Yard            string  `json:"yard"`
	ContainerNumber string  `json:"container_number"`
	Block           string  `json:"block"`
	Slot            int     `json:"slot"`
	Row             int     `json:"row"`
	Tier            int     `json:"tier"`
	ContainerSize   int     `json:"container_size"`
	ContainerHeight float64 `json:"container_height"`
	ContainerType   string  `json:"container_type"`
}

type PickupContainerRequest struct {
	Yard            string `json:"yard"`
	ContainerNumber string `json:"container_number"`
}

// Response
type SuggestContainerResponse struct {
	Yard  string `json:"yard"`
	Block string `json:"block"`
	Slot  int    `json:"slot"`
	Row   int    `json:"row"`
	Tier  int    `json:"tier"`
}
