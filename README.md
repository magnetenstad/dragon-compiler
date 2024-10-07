# dragon-compiler

Compiles to C

## Examples

```cpp
struct Color {
    r Int = 100
    g Int = 50
    b Int
}

struct House {
    street String = "Unknown street"
    streetNumber Int
    color Color
}

house = House(
    street "Kongens Gate"
    streetNumber 12
    color Color(
        r 254 + 1
    )
)

{
    skip_if house.streetNumber < 0

    print house.street
    print house.streetNumber
}
```

The example above compiles to the following C code.

```c
#include <stdio.h>

typedef struct {
    int r;
    int g;
    int b;
} Color;
void __Construct_Color__(Color *o) {
    o->r = 100;
    o->g = 50;
    o->b = 0;
}
typedef struct {
    char* street;
    int streetNumber;
    Color color;
} House;
void __Construct_House__(House *o) {
    o->street = "Unknown street";
    o->streetNumber = 0;
    __Construct_Color__(&o->color);
}

int main(int argc, char *argv[]) {;
    House __Instance_1__;
    __Construct_House__(&__Instance_1__);
    __Instance_1__.street = "Kongens Gate";
    __Instance_1__.streetNumber = 12;
    Color __Instance_2__;
    __Construct_Color__(&__Instance_2__);
    __Instance_2__.r = (254+1);
    __Instance_1__.color = __Instance_2__;
    House house = __Instance_1__;
    __StartBlock_1__: {;
        if ((house.streetNumber<0)) goto __EndBlock_1__;
        printf(house.street);
        printf(house.streetNumber);
    }
    __EndBlock_1__: {}
    return 0;
}
```
