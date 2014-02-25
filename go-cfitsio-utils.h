#ifndef GO_CFITSIO_UTILS_H
#define GO_CFITSIO_UTILS_H 1

#include <stdlib.h>
#include <string.h>

static
char* 
CStringN(int sz) 
{
 char* array = (char*)malloc(sz*sizeof(char));
 return array;
}

static
char**
char_array_new(int sz)
{
  int i = 0;
  char **array = (char**)malloc(sz*sizeof(char*));
  for (i = 0; i < sz; i++) {
    array[i] = NULL;
  }
  return array;
}

static
void
char_array_set(char** array, int idx, char *value)
{
  array[idx] = value;
}

static
long*
long_array_new(int sz)
{
  long *array = (long*)malloc(sz*sizeof(long));
  return array;
}

#endif /* !GO_CFITSIO_UTILS_H */
