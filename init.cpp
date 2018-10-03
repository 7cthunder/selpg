#include <cstdio>
#include <time.h>
#include <cstdlib>
using namespace std;

int random(int m) {
    return rand()%m;
}

int main(){
    srand(time(NULL));
    for (int i = 0; i < 200; i++) {
        printf("Line %3d\n", i);
        if (random(10) < 1) 
            printf("\f");
    }

}
