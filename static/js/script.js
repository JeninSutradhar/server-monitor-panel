document.addEventListener('DOMContentLoaded', () => {
  let contentArea = document.getElementById('content-area');
  let serviceDiv = document.getElementById("serviceDiv");
  let taskListDiv = document.getElementById("taskDiv");
  const taskDescriptionInput = document.getElementById('taskDesc');
  const taskScheduleInput = document.getElementById('taskSchedule');

  function redirectToDashboard() {
      window.location.href = "/";
  }
function handleFetchError(response, error, message) {
         console.error(`Fetch error ${message} `, error, response);
         if(error.message === 'Failed to fetch' || error.message.startsWith('Network')) {
             redirectToDashboard();
         }  else {

             alert( `Http error type: ${response ? response.status :  'Not found' } ${message} : ${error}` )
         }
     }
  function renderPageContent(element) {
      let target = "";
      if (element) {
          target = element.getAttribute('data-target');
      }

      if (target === "dashboard") {
          contentArea.innerHTML = `
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div class="bg-white shadow-md p-4 rounded-lg text-center">
                  <h2 class="text-xl font-semibold mb-2">CPU Usage</h2>
                  <p class="text-gray-700" id="cpuUsage">Loading...</p>
              </div>
              <div class="bg-white shadow-md p-4 rounded-lg text-center">
                  <h2 class="text-xl font-semibold mb-2">Memory Usage</h2>
                  <p class="text-gray-700" id="memoryUsage">Loading...</p>
              </div>
              <div class="bg-white shadow-md p-4 rounded-lg text-center">
                  <h2 class="text-xl font-semibold mb-2">Disk Usage</h2>
                  <p class="text-gray-700" id="diskUsage">Loading...</p>
              </div>
              <div class="bg-white shadow-md p-4 rounded-lg text-center">
                  <h2 class="text-xl font-semibold mb-2">Network Stats</h2>
                  <p class="text-gray-700">Sent: <span id="netSent">Loading...</span>, Received: <span id="netRecv">Loading...</span></p>
              </div>
              <div class="bg-white shadow-md p-4 rounded-lg text-center">
                  <h2 class="text-xl font-semibold mb-2">Load Average</h2>
                  <p class="text-gray-700">
                      <span id="load1">Loading...</span>,
                      <span id="load5">Loading...</span>,
                      <span id="load15">Loading...</span>
                  </p>
              </div>
              <div class="bg-white shadow-md p-4 rounded-lg text-center">
                  <h2 class="text-xl font-semibold mb-2">Uptime</h2>
                  <p class="text-gray-700" id="uptime">Loading...</p>
              </div>
              <div class="bg-white shadow-md p-4 rounded-lg text-center col-span-3">
                  <h2 class="text-xl font-semibold mb-2">Processes Running</h2>
                  <ul id="processList"></ul>
              </div>
              <div class="bg-white shadow-md p-4 rounded-lg text-center col-span-3">
                  <h2 class="text-xl font-semibold mb-2">Interfaces Info</h2>
                  <ul id="interfaceList"></ul>
              </div>
          </div>
          <div class="mt-8">
              <h2 class="text-xl font-semibold mb-4">Service Management</h2>
              <div id="serviceDiv" class="space-y-4"></div>
          </div>
          <div class="mt-8">
              <h2 class="text-xl font-semibold mb-4">Task Management</h2>
              <div class="space-y-4" id="taskDiv"></div>
              <div class="bg-white shadow-md p-4 rounded-lg space-y-2">
                  <h3 class="text-lg font-semibold mb-2">Schedule New Task:</h3>
                  <div class="flex">
                      <input type="text" class="border rounded px-2 mr-2" id="taskDesc" placeholder="Enter Description"/>
                      <input type="datetime-local" class="border px-2 rounded mr-2" id="taskSchedule"/>
                      <button class="bg-blue-500 text-white rounded hover:bg-blue-700 py-2 px-4 font-bold" onclick="handleNewTask()">Schedule</button>
                  </div>
              </div>
          </div>
          `;
          fetchData();
          handleServices();
          handleTasks();
      } else if (target === "firewall") {
          contentArea.innerHTML = `<div class="flex justify-center"> <h2 class="text-xl font-medium"> Page with firewall logic is under development</h2></div>`;
      } else if (target === "settings") {
          contentArea.innerHTML = `<div class="flex justify-center"> <h2 class="text-xl font-medium"> Page for settings is under development.</h2> </div>`;
      } else {
          redirectToDashboard();
      }
  }

  function fetchData() {
      fetch("/api/metrics")
          .then(response => {
              if (response.redirected) {
                  redirectToDashboard();
                  return;
              }
              if (!response.ok) {
                  throw new Error(`Error fetching metrics: ${response.status}`);
              }
              return response.json();
          })
          .then(data => {
              updateUI(data);
          })
          .catch(error => {
              handleFetchError(undefined, error, "fetching metrics");
          });
  }

  function handleServices() {
      fetch("/api/services", {
          method: "GET",
          headers: {
              'Content-Type': 'application/json'
          }
      })
          .then(response => {
              if (response.redirected) {
                  redirectToDashboard();
                  return;
              }
              if (!response.ok) {
                  throw new Error(`Error fetching services: ${response.status}`);
              }
              return response.json();
          })
          .then(data => {
              renderServices(data);
          })
          .catch(error => {
              handleFetchError(undefined, error, 'service info');
          });
  }

  function renderServices(services) {
      serviceDiv.innerHTML = "";
      for (const key in services) {
          const service = services[key];
          const div = document.createElement("div");
          div.classList.add("bg-white", "p-4", "rounded-lg", "shadow-md");

          const name = document.createElement("h3");
          name.classList.add("text-xl", "font-semibold", "mb-2");
          name.textContent = `Service: ${service.Name}`;
          div.appendChild(name);

          const status = document.createElement("p");
          status.classList.add("mb-2", "text-gray-700");
          status.textContent = `Status: ${service.Status}`;
          div.appendChild(status);

          const controls = document.createElement("div");
          controls.classList.add("flex", "space-x-2");

          const createButton = (text, className, action) => {
              const button = document.createElement("button");
              button.textContent = text;
              button.classList.add(...className.split(" "));
              button.onclick = () => handleServiceAction(service.Name, action);
              return button;
          };

          controls.appendChild(createButton("Install", "bg-blue-500 text-white rounded font-bold py-2 px-4 hover:bg-blue-700", "install"));
          controls.appendChild(createButton("Start", "bg-green-500 text-white py-2 px-4 rounded font-bold hover:bg-green-700", "start"));
          controls.appendChild(createButton("Stop", "bg-red-500 text-white py-2 px-4 rounded font-bold hover:bg-red-700", "stop"));
          controls.appendChild(createButton("Reload", "bg-yellow-500 text-white font-bold py-2 px-4 rounded hover:bg-yellow-700", "reload"));
          controls.appendChild(createButton("Uninstall", "bg-purple-500 text-white font-bold py-2 px-4 rounded hover:bg-purple-700", "uninstall"));


          div.appendChild(controls);
          serviceDiv.appendChild(div);
      }
  }

  function handleServiceAction(serviceName, action) {
      fetch("/api/services", {
          method: 'POST',
          headers: {
              "Content-Type": "application/json",
          },
          body: JSON.stringify({
              name: serviceName,
              action: action
          })
      })
          .then(response => {
              if (response.redirected) {
                  redirectToDashboard();
                  return;
              }
              if (!response.ok) {
                  throw new Error(`Error performing service action "${action}" for "${serviceName}": ${response.status}`);
              }
          })
          .then(() => {
              handleServices();
          })
          .catch(error => {
               handleFetchError(undefined, error, ' service action');
          });
  }

  function handleTasks() {
      fetch("/api/tasks", {
          method: "GET",
          headers: {
              "Content-Type": "application/json"
          }
      })
          .then(response => {
              if (response.redirected) {
                  redirectToDashboard();
                  return;
              }
              if (!response.ok) {
                  throw new Error(`Error fetching tasks: ${response.status}`);
              }
              return response.json();
          })
          .then(data => {
              renderTasks(data);
          })
          .catch(error => {
              handleFetchError(undefined, error , 'fetching tasks');
          });
  }

  function renderTasks(data) {
      taskListDiv.innerHTML = "";
      for (const id in data) {
          const task = data[id];
          const div = document.createElement("div");
          div.classList.add("bg-white", "shadow-md", "rounded-lg", "p-4", "flex", "items-center");

          const desc = document.createElement("span");
          desc.classList.add("font-medium", "mr-2");
          desc.textContent = `Description: ${task.Description}`;
          div.appendChild(desc);

          const status = document.createElement("span");
          status.classList.add("font-semibold", "text-gray-500");
          status.textContent = `Finished? ${task.IsFinished}`;
          div.appendChild(status);

          const removeBtn = document.createElement("button");
          removeBtn.classList.add("bg-red-500", "hover:bg-red-700", "font-bold", "py-1", "px-2", "rounded", "text-white", "ml-auto");
          removeBtn.textContent = "X";
          removeBtn.onclick = () => handleRemoveTask(task.ID);
          div.appendChild(removeBtn);
          taskListDiv.appendChild(div);
      }
  }

  function handleRemoveTask(taskId) {
      fetch(`/api/tasks?id=${taskId}`, {
          method: "DELETE",
      })
          .then(response => {
              if (response.redirected) {
                  redirectToDashboard();
                  return;
              }
              if (!response.ok) {
                  throw new Error(`Error deleting task ${taskId}: ${response.status}`);
              }
          })
          .then(() => {
              handleTasks();
          })
          .catch(error => {
             handleFetchError(undefined, error, `remove task ${taskId}` );

          });
  }

  function handleNewTask() {
      const description = taskDescriptionInput.value.trim();
      if (description === "") {
          alert("Please add a description before scheduling a new task!");
          return;
      }
      const timeValue = taskScheduleInput.value;
         if (timeValue === "") {
             alert("Time format is required!");
             return;
         }

      const runTime = new Date(timeValue).toJSON();

      fetch("/api/tasks", {
          method: 'POST',
          headers: {
              "Content-Type": "application/json"
          },
          body: JSON.stringify({
              description: description,
              run_time: runTime
          })
      })
          .then(response => {
              if (response.redirected) {
                  redirectToDashboard();
                  return;
              }
               if (!response.ok) {
                  throw new Error(`Error scheduling new task: ${response.status}`);
              }
                taskDescriptionInput.value = "";
                taskScheduleInput.value = "";

          })
          .then(() => {
              handleTasks();
          })
          .catch(error => {
               handleFetchError(undefined, error, 'new task ');

          });
  }

  function updateUI(data) {
      document.getElementById("cpuUsage").textContent = `${data.cpu.usage.toFixed(2)}%`;
      document.getElementById("memoryUsage").textContent = `${data.memory.used_percent.toFixed(2)}%`;
      document.getElementById("diskUsage").textContent = `${data.disk.used_percent.toFixed(2)}%`;
      document.getElementById("netSent").textContent = formatBytes(data.network.bytes_sent);
      document.getElementById("netRecv").textContent = formatBytes(data.network.bytes_recv);
      document.getElementById("load1").textContent = data.load.load1.toFixed(2);
      document.getElementById("load5").textContent = data.load.load5.toFixed(2);
      document.getElementById("load15").textContent = data.load.load15.toFixed(2);
      document.getElementById("uptime").textContent = formatUptime(data.host.uptime);

      const processList = document.getElementById("processList");
      processList.innerHTML = "";
      data.processes.forEach(process => {
          const li = document.createElement("li");
          li.classList.add("text-sm", "p-1", "border-b", "flex", "items-center", "space-x-2");
          const pid = document.createElement("span");
          pid.classList.add("text-gray-500");
          pid.textContent = `PID: ${process.Pid}`;
          const userName = document.createElement("span");
           userName.textContent = `User: ${process.Username}`;
          const processInfo = document.createElement("span");
          processInfo.textContent = `${process.Name} (CPU ${process.CPUUsage.toFixed(2)}%) (Mem: ${process.MemUsage.toFixed(2)}%)`;
          li.appendChild(pid);
          li.appendChild(userName);
           li.appendChild(processInfo);
          processList.appendChild(li);
      });

      const interfacesList = document.getElementById("interfaceList");
      interfacesList.innerHTML = "";
      data.network_interfaces.interfaces.forEach(item => {
          const list = document.createElement("li");
          list.classList.add("bg-gray-100", "mb-1", "shadow-sm", "p-2", "rounded", "border");

          const nameInterface = document.createElement("h4");
          nameInterface.classList.add("text-base", "font-medium", "mb-2");
          nameInterface.textContent = `Interface: ${item.Name}`;
          list.appendChild(nameInterface);

          const listIp = document.createElement("ul");
          item.IPs.forEach(ip => {
              const li = document.createElement("li");
              li.textContent = ip;
              listIp.appendChild(li);
          });
          const mac = document.createElement("p");
          mac.textContent = `Mac Address: ${item.MacAddress}`;
          list.appendChild(mac);
          list.appendChild(listIp);
          interfacesList.appendChild(list);
      });
  }

  function formatBytes(bytes) {
      if (bytes === 0) return "0 Bytes";
      const k = 1024;
      const sizes = ["Bytes", "KB", "MB", "GB", "TB"];
      const i = parseInt(Math.floor(Math.log(bytes) / Math.log(k)));
      return `${parseFloat((bytes / Math.pow(k, i)).toFixed(2))} ${sizes[i]}`;
  }

  function formatUptime(totalSeconds) {
      const days = Math.floor(totalSeconds / (60 * 60 * 24));
      const hours = Math.floor((totalSeconds % (60 * 60 * 24)) / (60 * 60));
      const minutes = Math.floor((totalSeconds % (60 * 60)) / 60);
      const seconds = Math.floor(totalSeconds % 60);
      return `${days}d ${hours}h ${minutes}m ${seconds}s`;
  }
  renderPageContent();
  fetchData();
  handleServices();
  handleTasks();
  setInterval(() => {
      fetchData();
  }, 5000);
  setInterval(() => {
      handleServices();
  }, 7000);
  setInterval(() => {
      handleTasks();
  }, 6000);
});