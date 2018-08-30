#include <stdio.h>
#include <string.h>
#include <stdlib.h>

#define ROOMX 12
#define ROOMY 12
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

void iteration(int *room, int *temp){
   int x, y;
   for(y = 0; y < ROOMY; y++){
      for(x = 0; x < ROOMX; x++){
         temp[y * ROOMX + x] = newValue(room, y, x);
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
   double k = 1;
   int new = old[getPos(row, ROOMY) * ROOMX + getPos(col, ROOMX)];
   new += k*old[getPos(row+1, ROOMY) * ROOMX + getPos(col, ROOMX)];
   new += k*old[getPos(row-1, ROOMY) * ROOMX + getPos(col, ROOMX)];
   new += k*old[getPos(row, ROOMY) * ROOMX + getPos(col+1, ROOMX)];
   new += k*old[getPos(row, ROOMY) * ROOMX + getPos(col-1, ROOMX)];
   //new -= k*(4*old[getPos(row, ROOMY) * ROOMX + getPos(col, ROOMX)]);
   
   new = new / 5;

   return new;
}


int main(int argc, char **argv)
{
   int t;

   //set initial room temperatures
   int *room, *temp, *temp2;
   room = malloc(ROOMX*ROOMY * sizeof(int));
   temp = malloc(ROOMX*ROOMY * sizeof(int));
   memset(room, 0, ROOMX * ROOMY * sizeof(int));
   room[3 * ROOMX + 3] = 80;
   room[10 * ROOMX + 10] = 80;
   printArray(room, ROOMX, ROOMY);

   //iterate over time
   for(t = 0; t < MAX_TIME; t++){
      
      printArray(room, ROOMX, ROOMY);
      
      //serial implementation
      iteration(room, temp);
      
      //swap arrays
      temp2 = room;
      room = temp;
      temp = temp2;
      
      room[3 * ROOMX + 3] = 80;
      room[10 * ROOMX + 10] = 80;
   }
   
   printArray(room, ROOMX, ROOMY);
   
}