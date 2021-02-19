package ses

import (
	"bytes"
	"crypto/md5"
	"fmt"

	"github.com/rs/xid"
)

func SessionID() string {
	sb := new(bytes.Buffer)
	sb.WriteString(xid.New().String())
	return fmt.Sprintf("%x", md5.Sum(sb.Bytes()))
}
