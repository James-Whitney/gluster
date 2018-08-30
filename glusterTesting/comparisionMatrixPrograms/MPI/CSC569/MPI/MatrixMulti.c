#include <stdio.h>
#include "mpi.h"


#define N 2048

/* global */
int a[N*N];

int b[N*N];

int c[N*N];

int fillArrays()
{
   int i, j;
   for(i = 0; i < N; i++){
      for(j = 0; j < N; j++){
         a[i * N + j] = 1;
         b[i * N + j] = 2;
      }
   }
}

int main(int argc, char *argv[])
{
  MPI_Status status;
  int me,p;
  int i,j;

  fillArrays();

  /* Start up MPI */

  MPI_Init(&argc, &argv);
  MPI_Comm_rank(MPI_COMM_WORLD, &me);
  MPI_Comm_size(MPI_COMM_WORLD, &p);

  printf("me=%d, p=%d", me, p);
   int row, col, pValue, k;
   for(row = start; row < end; row++) {
      for(col = 0; col < width; col++) {
         pValue = 0;
         for(k = 0; k < width; k++) {
            pValue += a[row * width + k] * b[k * width + col];
         }
         c[row * width + col] = pValue;
      }
   }
  /* Result gathering */
  if (me != 0 )
  {
      MPI_Send(&c[(me) * N/p][0], N*N/p, MPI_INT, 0, 0, MPI_COMM_WORLD);
  }
  else
  {
      for (i=1; i<p; i++)
      {
          MPI_Recv(&c[i * N/p][0], N * N / p, MPI_INT, i, 0, MPI_COMM_WORLD, 0);
      }
  }

  MPI_Barrier(MPI_COMM_WORLD);


  /* Quit */

  MPI_Finalize();
  return 0;
}