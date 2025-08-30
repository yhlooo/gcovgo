#include <stdio.h>

#include "path/to/lib1/lib1.h"
#include "branches/branches.h"

int main() {
    printf("Hello World!\n");

    for (int i = 0; i < 13; i++) {
        printf("switch case %d: %d\n", i, sample_switch(i));
    }

    printf("if 3: %d\n", sample_if(3));
    printf("if 8: %d\n", sample_if(8));
    printf("if 20: %d\n", sample_if(20));
    printf("if 24: %d\n", sample_if(24));
    printf("if 25: %d\n", sample_if(25));

    printf("loop 233: %d\n", sample_loop(233));

    func1();
    func2(66, 233);

    return 0;
}
