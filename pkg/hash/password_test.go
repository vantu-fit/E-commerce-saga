package hash

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComparePassword(t *testing.T){
	password :=  "deptrai123"
	hashedPassword, err := HashedPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	fmt.Println("Hashed Password: ", hashedPassword)
	fmt.Println("Password: ", password)

	err = CheckPassword(hashedPassword, password)
	require.NoError(t, err)
}