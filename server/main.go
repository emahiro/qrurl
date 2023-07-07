package main

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	qrurlv1 "github.com/emahiro/qrurl/server/gen/proto/qrurl/v1"
)

type QrUrlService struct{}

func (s *QrUrlService) PostCode(
	ctx context.Context,
	req *connect.Request[qrurlv1.PostCodeRequest],
) {
}

func main() {
	fmt.Println("Hello, World!")
}
