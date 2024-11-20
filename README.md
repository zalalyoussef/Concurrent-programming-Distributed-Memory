## **Distributed Memory Program (L2)**: **Data Processing with MPI or Go Channels**

### **Overview**

In this task, we modify the previous L1 program so that no shared memory exists, and all communication between processes is done via messages or channels. This means the data and result managers will handle the data and results independently, and the main process will handle communication through requests. The workers will process data in parallel and send results back to the result manager. 

This task utilizes **distributed memory** where data and results are passed between separate processes, avoiding shared memory structures. The goal is to implement message-based communication using either **MPI** (in C++) or **Go Channels**.

### **Program Structure**

There are four key processes:

1. **Main Process**: Responsible for reading data, spawning worker, data manager, and result manager processes, sending data to the data manager, and receiving results from the result manager.
   
2. **Data Process**: Holds an internal array and accepts requests to insert or remove items. It cannot send messages on its own but can only process requests. The data array size is limited to half the number of elements in the data file.
   
3. **Worker Processes**: Receive data items from the data process, compute results, filter them based on a criterion, and send valid results to the result process.
   
4. **Result Process**: Receives valid results from worker processes, stores them in an internal array, and when all work is done, sends all items back to the main process.

### **Flow of Communication**

1. **Main Process**:
   - Reads data from a file.
   - Spawns the worker processes, a data manager, and a result manager.
   - Sends the data to the data manager process one by one.
   - Receives the final results from the result process.
   - Outputs the sorted results to a file.

2. **Data Process**:
   - Stores the data in an array.
   - Receives requests to insert or remove items from the array.
   - It can only accept insertion requests when there is space, and only removal requests when there is data in the array.
   
3. **Worker Processes**:
   - Request data items from the data process.
   - Process each item using a selected function.
   - If the result meets a filter criterion, send it to the result process.
   
4. **Result Process**:
   - Stores the received results in a sorted array.
   - Sends the final sorted results back to the main process when all items are processed.

### **Key Concepts**

- **No Shared Memory**: The program avoids shared memory, which means there are no global data structures accessible by multiple processes. All data is passed as messages or through channels.

- **Message Passing**: Communication between processes is done through **messages** in MPI (C++) or **channels** in Go. The processes send and receive data requests and results through these mechanisms.

- **Message Types**: 
  - *Insert Item*: Main or worker processes send insert requests to the data manager.
  - *Remove Item*: Worker processes request items from the data manager.
  - *Result*: Worker processes send computed results to the result manager.

- **Filter Criterion**: The results that workers send to the result process must meet a specific filter criterion, which can be based on computed values (e.g., whether the result is even, greater than a certain threshold, etc.).

- **Data and Result Managers**: 
  - *Data Manager*: Manages the data array and ensures synchronization when inserting or removing items.
  - *Result Manager*: Manages the result array and ensures that received results are stored in sorted order.

### **Required Tools**

You need to choose one of the following toolsets for implementing the distributed memory solution:

- **Go** (with channels for synchronization)
- **C++ with MPI** (using MPI processes and communicators)

#### **Installation Requirements:**

**For MPI (C++):**
- **Linux**: 
  - Install g++ and OpenMPI:
    ```bash
    sudo apt-get install g++ libopenmpi-dev openmpi-bin
    ```
  
**For Go:**
- **Install Go** (if not already installed):
  ```bash
  wget https://golang.org/dl/go1.18.1.linux-amd64.tar.gz
  sudo tar -C /usr/local -xvzf go1.18.1.linux-amd64.tar.gz
  ```

- **Use Go Channels**: Channels are used to handle the communication between processes, sending and receiving data, and synchronizing the actions of workers, data managers, and result managers.

### **Steps to Implement**

1. **Main Process**:
   - Read the input data file and store it in an array or list.
   - Spawn the worker, data manager, and result manager processes.
   - Send data items to the data manager.
   - Collect results from the result manager.
   - Output the results to a file.

2. **Data Process**:
   - Store the data in a private array.
   - Handle insert and remove requests from the main process or workers.
   - Ensure synchronization when removing or inserting items based on array capacity.

3. **Worker Processes**:
   - Request items from the data manager.
   - Perform calculations on the received item.
   - Filter results and send them to the result manager if they match the criterion.

4. **Result Process**:
   - Store the results received from the workers in a sorted array.
   - Send all processed items back to the main process.

### **File Formats**

- **Input Data**: Same as in Program L1 — JSON formatted file containing student data.
- **Result File**: A text file containing the results, sorted by hash, and including the computed sums of relevant fields (e.g., `int` and `float`).

Example:
```json
{
  "students": [
    {"name": "Antanas", "year": 1, "grade": 6.95},
    {"name": "Kazys", "year": 2, "grade": 8.65},
    {"name": "Petras", "year": 2, "grade": 7.01}
    // More students...
  ]
}
```

Output:
```txt
Name | Year | Grade | Hash
----------------------------------------
Jonas | 1 | 6.95 | A18EAC8F30AC0FC630AE175A851CA5DA24FA8C85
Kazys | 2 | 8.65 | E763185F4B7303A787ACC513B7AA56706C7A42AC
// More results...
```

### **Usage**

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/DistributedMemoryDataProcessing.git
   ```

2. **For MPI**:
   - Compile with MPI:
     ```bash
     mpic++ -o program program.cpp -fopenmp
     ```
   - Run the program:
     ```bash
     mpiexec -np <number_of_processes> ./program
     ```

3. **For Go**:
   - Build the Go project:
     ```bash
     go build program.go
     ```
   - Run the program:
     ```bash
     ./program
     ```

4. The program will read data from the input file, process it with worker processes, and output the results to the result file.

---
