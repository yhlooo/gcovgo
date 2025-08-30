// switch 样例
int sample_switch(int in) {
    int ret = 0;
    switch (in) {
    case 1:
        ret = 1;
        break;
    case 2:
        ret = 2;
        break;
    case 10:
        // += 10
        ret--;
    case 11:
        // += 11
        ret--;
    case 12:
        ret += 12;
        break;
    }

    return ret;
}

// if 样例
int sample_if(int in) {
    if (in > 10) {
        if (in % 10 == 0 || in % 3 == 0) {
            return 11;
        }
        return 10;
    } else if (in > 5) {
        int ret = 20 + (in % 10);
        if (in > 7 && in % 2 == 0) {
            ret -= 5;
        } else {
            ret -= 5;
            ret += 7;
        }
        return ret;
    } else {
        // empty
    }

    return 30;
}

// 循环样例
int sample_loop(int in) {
    for (int i = 0; i < 10; i++) {
        in += i;
        in %= 10;
        for (int j = 0; j < 10; j++) {
            in += j * 10;
            if (j == 5) {
                continue;
            }
        }
    }

    int i = 10;
    do {
        i--;
        in -= i;
    } while (in > 0 || i > 0);

    while (i < 100) {
        i++;
        in %= i;
        if (in == 0) {
            break;
        }
    }

    return in;
}
