#include <stdio.h>

void func1() {
    printf("func1 called\n");
}

int func2(int i, int j) {
    printf("func2(%d, %d) called\n", i, j);
    return i + j;
}
