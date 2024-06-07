# Documentation

!!! IMPORTANT !!! If you want to contribute to the documentation, please read the [Contributing](https://github.com/pewpewlive/hybroid/blob/master/CONTRIBUTING.md) file.

The documentation of the codebase is spread throughout it. With each folder from root comes a README.md, where the doc lies, covering everything inside that folder. 

There is some information that will not exist in any of those docs though (simply because it wouldn't make sense). This file will provide such information.

## Struct implementations of interfaces

When implementating a struct, which you want to implement a certain interface, the supposed methods that exist to adhere to the interface must not be tailored to the type (T), but to the pointer of the type (*T).

There are a few reasons for this. First of all it enforces us to use pointers to share data, and not values themselves, preventing copying and making for better memory efficiency. 