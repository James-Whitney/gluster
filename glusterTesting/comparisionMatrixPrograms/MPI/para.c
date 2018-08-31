#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <mpi.h>

#define GRID_X 12
#define GRID_Y 12
#define MAX_TIME 20

void printArray(int *arr, int width, int height){
   int x, y;
   for(y = 0; y < height; y++){
      for(x = 0; x < width; x++){
         if(arr[y*width + x] == 0)
            printf("---- ");
         else
            printf("%.4d ", arr[y*width + x]);
      }
      printf("\n");
   }
   printf("\n");
}

void iteration(int *grid, int *temp, int my_rank, int comm_sz){
   int col, row;
   int start = my_rank * GRID_Y / comm_sz;
   int end = (my_rank+1) * GRID_Y / comm_sz;
   for(row = start; row < end; row++){
      for(col = 0; col < GRID_X; col++){
         temp[row * GRID_Y + col] = newValue(grid, row, col);
      }
   }
}

int getPos(int pos, int dim){
   if(pos < 0)
      return 0;
   if(pos >= dim)
      return dim-1;
   return pos;
}

int newValue(int *old, int row, int col)
{
   int new = old[getPos(row, GRID_Y) * GRID_X + getPos(col, GRID_X)];
   new += old[getPos(row+1, GRID_Y) * GRID_X + getPos(col, GRID_X)];
   new += old[getPos(row-1, GRID_Y) * GRID_X + getPos(col, GRID_X)];
   new += old[getPos(row, GRID_Y) * GRID_X + getPos(col-1, GRID_X)];
   new += old[getPos(row, GRID_Y) * GRID_X + getPos(col+1, GRID_X)];
   new = new / 5;
   return new;
}


int main(int argc, char **argv)
{
   int my_rank, comm_sz;
   double startTime, endTime;
   startTime = MPI_Wtime();
   int t;

   MPI_Init(NULL, NULL);
   MPI_Comm_rank(MPI_COMM_WORLD, &my_rank);
   MPI_Comm_size(MPI_COMM_WORLD, &comm_sz);

   //set initial grid temperatures
   int *grid, *temp, *swap;
   grid = malloc(GRID_X*GRID_Y * sizeof(int));
   temp = malloc(GRID_X*GRID_Y * sizeof(int));
   memset(grid, 0, GRID_X * GRID_Y * sizeof(int));
   grid[3 * GRID_X + 3] = 80;
   grid[10 * GRID_X + 10] = 80;
   if (my_rank == 0) {
      MPI_Bcast(grid, GRID_X*GRID_Y, MPI_INT, 0, MPI_COMM_WORLD);
   }
   else {
      MPI_Recv(grid, GRID_X*GRID_Y, MPI_INT, 0, 0, MPI_COMM_WORLD, MPI_STATUS_IGNORE);
   }

   //iterate over time
   for(t = 0; t < MAX_TIME; t++){

      printArray(grid, GRID_X, GRID_Y);      

      iteration(grid, temp, my_rank, comm_sz);
      
      MPI_Allgather(temp, GRID_X*GRID_Y, MPI_INT, grid, GRID_X*GRID_Y, MPI_INT, MPI_COMM_WORLD);

      grid[3 * GRID_X + 3] = 80;
      grid[10 * GRID_X + 10] = 80;
   }
   if (my_rank == 0){
      printArray(grid, GRID_X, GRID_Y);
      endTime = MPI_Wtime();
      printf("Run Time: %f\n", endTime - startTime);
   }
   
   MPI_Finalize();
}
