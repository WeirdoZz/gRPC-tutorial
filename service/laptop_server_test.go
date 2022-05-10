package service_test

import (
	"context"
	"gRPC-tutori/pb"
	"gRPC-tutori/sample"
	"gRPC-tutori/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestLaptopServer_CreateLaptop(t *testing.T) {
	t.Parallel()

	laptopNoID := sample.NewLaptop()
	laptopNoID.Id = ""
	laptopInvalidID := sample.NewLaptop()
	laptopInvalidID.Id = "invalid-uuid"
	laptopDuplicateID := sample.NewLaptop()
	storeDuplicateID := service.NewInMemoryLaptopStore()
	err := storeDuplicateID.Save(laptopDuplicateID)
	require.Nil(t, err)

	testCases := []struct {
		name        string
		laptop      *pb.Laptop
		laptopStore service.LaptopStore
		imageStore  service.ImageStore
		ratingStore service.RatingStore
		code        codes.Code
	}{
		{
			name:        "success_with_id",
			laptop:      sample.NewLaptop(),
			laptopStore: service.NewInMemoryLaptopStore(),
			imageStore:  service.NewDiskImageStore("img"),
			ratingStore: service.NewInMemoryRatingStore(),
			code:        codes.OK,
		},
		{
			name:        "success_no_id",
			laptop:      laptopNoID,
			laptopStore: service.NewInMemoryLaptopStore(),
			imageStore:  service.NewDiskImageStore("img"),
			ratingStore: service.NewInMemoryRatingStore(),
			code:        codes.OK,
		},
		{
			name:        "failure_invalid_id",
			laptop:      laptopInvalidID,
			laptopStore: service.NewInMemoryLaptopStore(),
			imageStore:  service.NewDiskImageStore("img"),
			ratingStore: service.NewInMemoryRatingStore(),
			code:        codes.InvalidArgument,
		},
		{
			name:        "failure_duplicate_id",
			laptop:      laptopDuplicateID,
			laptopStore: storeDuplicateID,
			imageStore:  service.NewDiskImageStore("img"),
			ratingStore: service.NewInMemoryRatingStore(),
			code:        codes.AlreadyExists,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := &pb.CreateLaptopRequest{Laptop: tc.laptop}

			server := service.NewLaptopServer(tc.laptopStore, tc.imageStore, tc.ratingStore)
			res, err := server.CreateLaptop(context.Background(), req)

			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.Id)
				if len(tc.laptop.Id) > 0 {
					require.Equal(t, tc.laptop.Id, res.Id)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, tc.code, st.Code())
			}
		})
	}
}
