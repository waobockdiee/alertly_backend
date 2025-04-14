FROM ubuntu:20.04 AS opencv-builder

# Instala dependencias necesarias
RUN apt-get update && apt-get install -y \
    build-essential cmake git pkg-config \
    libgtk-3-dev libavcodec-dev libavformat-dev libswscale-dev \
    libv4l-dev libxvidcore-dev libx264-dev libjpeg-dev \
    libpng-dev libtiff-dev gfortran openexr libatlas-base-dev \
    python3-dev python3-numpy libtbb2 libtbb-dev libdc1394-22-dev

# Descarga OpenCV 4.11.0 y opencv_contrib 4.11.0
WORKDIR /opt
RUN git clone --branch 4.11.0 https://github.com/opencv/opencv.git
RUN git clone --branch 4.11.0 https://github.com/opencv/opencv_contrib.git

# Compila OpenCV
WORKDIR /opt/opencv/build
RUN cmake -D CMAKE_BUILD_TYPE=Release \
          -D CMAKE_INSTALL_PREFIX=/usr/local \
          -D OPENCV_EXTRA_MODULES_PATH=/opt/opencv_contrib/modules \
          -D BUILD_EXAMPLES=OFF ..
RUN make -j$(nproc) && make install

# En la siguiente etapa, puedes copiar los archivos instalados a la etapa builder de tu aplicación
