#ifndef GO_CFITSIO_UTILS_H
#define GO_CFITSIO_UTILS_H 1

#include <stdlib.h>
#include <string.h>

static
char* 
char_buf_array(int sz) 
{
 char* array = (char*)malloc(sz*sizeof(char));
 return array;
}

#endif /* !GO_CFITSIO_UTILS_H */
