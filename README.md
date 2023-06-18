# Fractals

An experimental project for my scientific hobbit, I am going to work with fractals of the Mandelbrot set (maybe some others) and image management with golang

## Use
```bash
git clone https://github.com/karchx/fractals.git
cd fractals
go build 
./fractals -h  # to see list of available customizations
./fractals -height 1000 -width 1000
```

## About the Algorithm
### The Math in a Nutshell
The Mandelbrot set is defined as the set of complex numbers $z_0$ for which the series 

$$z_{n+1} = z²_n + z_0$$

is bounded for all $n ≥ 0$. In other words, $z_0$ is part of the Mandelbrot set if $z_n$ does not approach infinity. This is equivalent to the  magnitude $|z_n| ≤ 2$ for all $n ≥ 0$.
