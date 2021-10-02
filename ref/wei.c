#include <stdio.h>
#include <stdint.h>
#include <unistd.h>
/*
void print(uint64_t value) {
    char buffer[65];
    size_t count = 0;

    buffer[65 - count++ - 1] = '\n';
    do {
        buffer[65 - count++ - 1] = (value % 10) + '0';
        value /= 10;
    } while (value > 0);
    buffer[65 - count++ - 1] = '\0';

    fwrite(&buffer[65 - count], 1, count, stdout);
}
*/
void print(uint64_t value) {
    char buffer[65];
    uint32_t count = 0;

    buffer[65 - count++ - 1] = '\n';
    do {
        buffer[65 - count++ - 1] = (value % 10) + '0';
        value /= 10;
    } while (value > 0);
    buffer[65 - count++ - 1] = '\0';

    write(1, &buffer[65 - count], count);
}

int main() {
    int lol[2];
    lol[0] = 34;
    lol[1] = 35;
    lol[0] = lol[0] + lol[1];
    print(lol[0]);
    return lol[0];
}

