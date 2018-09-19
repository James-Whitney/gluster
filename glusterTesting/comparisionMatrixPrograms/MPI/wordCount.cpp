#include <string>
#include <cstdlib>
#include <cerrno>
#include <iostream>
#include <fstream>
#include <map>
#include <algorithm>
#include <functional>
#include <cctype>
#include <vector>
#include <sstream>
#include <iterator>
#include <omp.h>
#include <ctime>
#include "mpi.h"

#define NUM_THREADS 28

int main(int argc, char *argv[])
{
    std::clock_t start_clock = std::clock();
    /* Start up MPI */
    MPI_Status status;
    int me,p;
    MPI_Init(&argc, &argv);
    MPI_Comm_rank(MPI_COMM_WORLD, &me);
    MPI_Comm_size(MPI_COMM_WORLD, &p);
    printf("me=%d, p=%d", me, p);

    // Node 0 duties:
    if (me == 0) {
        std::string book;

        // Read in the book into a C++ string
        std::FILE *fp = std::fopen("words.txt", "r");
        if (fp) {
            std::fseek(fp, 0, SEEK_END);
            book.resize(std::ftell(fp));
            std::rewind(fp);
            std::fread(&book[0], 1, book.size(), fp);
            std::fclose(fp);
        }
        else {
            throw(errno);
        }

        //repeat book x times
        std::string temp = book;
        for (int i = 1; i < 10; i++) {
            book += temp;
        }

        //process the book into a vector of words
        std::stringstream ss(book);
        std::istream_iterator<std::string> begin(ss);
        std::istream_iterator<std::string> s_end;
        std::vector<std::string> words(begin, s_end);

        // int start = me * words.size() / (p - 1);
        // int end = (me + 1) *  words.size() / (p - 1);

        //send chunk of book to workers
        int start, end;
        for (int i = 1; i < p; i ++) {
            start = i * words.size() / (p - 1);
            end = (i + 1) *  words.size() / (p - 1);
            std::string wordChunk = words[start];
            for (int j = start + 1; i < end; j++) {
                wordChunk += (words[j] + " ");
            }
            int size = wordChunk.length() + 1;
            MPI_Send(&size, 1, MPI_INT, i, 0, MPI_COMM_WORLD);
            MPI_Send(&wordChunk.c_str(), wordChunk.length() + 1, MPI_CHAR, i, 0, MPI_COMM_WORLD);
        }

        int recvCountSize;
        int recvListSize;
        char * wordSubList;
        int * wordSubCount;
        std::map<std::string,int> m;
        omp_lock_t writelock;
        omp_init_lock(&writelock);
        for (int i = 1; i < p; i ++) {
            MPI_Recv(recvListSize, 1, MPI_INT, i, 0, MPI_COMM_WORLD, 0);
            wordSubList = (char*)malloc(sizeof(char) * recvListSize);
            if (wordSubList == NULL) {
                return 1;
            }
            MPI_Recv(wordSubList, recvListSize, MPI_CHAR, i, 0, MPI_COMM_WORLD, 0);
            MPI_Recv(recvCountSize, 1, MPI_INT, i, 0, MPI_COMM_WORLD, 0);
            wordSubCount = (int*)malloc(sizeof(int) * recvCountSize);
            if (wordSubCount == NULL) {
                return 1;
            }
            MPI_Recv(wordSubCount, recvCountSize, MPI_INT, i, 0, MPI_COMM_WORLD, 0);

            std::string subBook = wordSubList;
            std::stringstream ss(subBook);
            std::istream_iterator<std::string> begin(ss);
            std::istream_iterator<std::string> s_end;
            std::vector<std::string> words(begin, s_end);

            int count;
            #pragma omp parallel for private(count, m, writelock) num_threads(NUM_THREADS)
            for (int j=0; j < recvCountSize; j++) {
                count = m[words[j]];
                if (count == 0) {
                    omp_set_lock(&writelock);
                    m[words[j]] = wordSubCount[j];
                    omp_unset_lock(&writelock);
                }
                else {
                    omp_set_lock(&writelock);
                    m[words[j]] = count + wordSubCount[j];
                    omp_unset_lock(&writelock);
                }
            }

            free(wordSubList);
            free(wordSubCount);

            double duration = ( std::clock() - start_clock ) / (double) CLOCKS_PER_SEC;
            std::cout << "time:" << duration << std::endl;
        }

    }
    else {
        int size;
        MPI_Recv(&size, 1, MPI_INT, 0, 0, MPI_COMM_WORLD, 0);
        char * wordChunk;
        wordChunk = (char*) malloc(size);
        if (wordChunk == NULL) {
            return 1;
        }
        MPI_Recv(wordChunk, size, MPI_CHAR, 0, 0, MPI_COMM_WORLD, 0);
        std::string book = wordChunk;

        //reprocess book into words again
        std::stringstream ss(book);
        std::istream_iterator<std::string> begin(ss);
        std::istream_iterator<std::string> end;
        std::vector<std::string> words(begin, end);

        std::map<std::string,int> results;
        omp_lock_t writelock;
        omp_init_lock(&writelock);
        int count;
        #pragma omp parallel for private(count) num_threads(NUM_THREADS)
        for (unsigned i=0; i < words.size(); i++) {
            count = results[words[i]];
            if (count == 0) {
                omp_set_lock(&writelock);
                results[words[i]] = 1;
                omp_unset_lock(&writelock);
            }
            else {
                omp_set_lock(&writelock);
                results[words[i]] = count + 1;
                omp_unset_lock(&writelock);
            }
        }
        omp_destroy_lock(&writelock);

        std::string wordList = "";
        std::vector<int> wordCounts;
        for (std::map<std::string,int>::iterator it=results.begin(); it!=results.end(); ++it) {
            wordList += (it->first + " ");
            wordCounts.push_back(it->second);
        }
        int returnListSize = wordList.length() + 1;
        int returnCountSize = results.size();

        MPI_Send(&returnListSize, 1, MPI_INT, 0, 0, MPI_COMM_WORLD);
        MPI_Send(&wordList.c_str(), returnListSize, MPI_CHAR, 0, 0, MPI_COMM_WORLD);
        MPI_Send(&returnCountSize, 1, MPI_INT, 0, 0, MPI_COMM_WORLD);
        MPI_Send(&wordCount[0], returnCountSize, MPI_INT, 0, 0, MPI_COMM_WORLD);
        free(wordChunk);
    }

    return 0;
}