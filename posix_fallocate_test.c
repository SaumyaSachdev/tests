#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <sys/time.h>

int main() {
    const char *filename = "testfile.txt";
    off_t file_size = 8 * 1024 * 1024; // 8 MB
    int fd;
    struct timeval start_time, end_time;
    double elapsed_time;

    // Open the file
    fd = open(filename, O_WRONLY | O_CREAT | O_EXCL , S_IRUSR | S_IWUSR);
    if (fd == -1) {
        perror("open");
        return 1;
    }

    // Measure start time
    gettimeofday(&start_time, NULL);

    // Allocate space for the file
    if (posix_fallocate(fd, 0, file_size) != 0) {
    

    // if (fallocate(fd, 0, 0, file_size) != 0) {
        perror("posix_fallocate");
        close(fd);
        return 1;
    }


    close(fd);
    fd = open(filename, O_WRONLY | O_CREAT  | O_DSYNC, S_IRUSR | S_IWUSR);
    if (fd == -1) {
        perror("open");
        return 1;
    }

    const char data[8 * 1024] = "This is 8KB of data.\n";
    ssize_t bytes_written = write(fd, data, sizeof(data));
    if (bytes_written == -1) {
        perror("write");
        close(fd);
        return 1;
    }

    // Measure end time
    gettimeofday(&end_time, NULL);

    
    // Calculate elapsed time in milliseconds
    elapsed_time = (end_time.tv_sec - start_time.tv_sec) * 1000.0; // seconds to milliseconds
    elapsed_time += (end_time.tv_usec - start_time.tv_usec) / 1000.0; // microseconds to milliseconds

    printf("Time taken to allocate %ld bytes: %.2f milliseconds\n", (long)file_size, elapsed_time);

    // Close and delete the file
    close(fd);
    // if (unlink(filename) == -1) {
    //     perror("unlink");
    //     return 1;
    // }

    return 0;
}
