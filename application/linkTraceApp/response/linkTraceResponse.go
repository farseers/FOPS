package response

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/trace"
)

type LinkTraceResponse struct {
	Entry trace.TraceContext
	List  collections.List[LinkTraceVO]
}
