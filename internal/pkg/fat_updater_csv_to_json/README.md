docker build -t mi-compilador-pyinstaller .
docker run --rm -v "$(pwd):/app/dist" mi-compilador-pyinstaller
