document.addEventListener("DOMContentLoaded", () => {
    // Fetch dashboard summary data
    fetchDashboardSummary();

    // Initialize charts
    initializeFolderChart();
    initializeUserActivityChart();
});

/**
 * Fetch and display dashboard summary data (total folders, files, users)
 */
function fetchDashboardSummary() {
    fetch("/api/dashboard/summary")
        .then((response) => {
            if (!response.ok) {
                throw new Error("Failed to fetch dashboard summary data.");
            }
            return response.json();
        })
        .then((data) => {
            document.getElementById("total-folders").innerText = data.totalFolders || 0;
            document.getElementById("total-files").innerText = data.totalFiles || 0;
            document.getElementById("total-users").innerText = data.totalUsers || 0;
        })
        .catch((error) => {
            console.error("Error fetching dashboard summary:", error);
        });
}

/**
 * Initialize the Folder Chart (Bar Chart)
 */
function initializeFolderChart() {
    const ctx = document.getElementById("folderChart").getContext("2d");

    // Simulated folder data for demonstration (replace with API data if needed)
    const data = [10, 20, 15, 30];
    const labels = ["Jan", "Feb", "Mar", "Apr"];

    new Chart(ctx, {
        type: "bar",
        data: {
            labels: labels,
            datasets: [
                {
                    label: "Folders Created",
                    data: data,
                    backgroundColor: "rgba(54, 162, 235, 0.6)",
                    borderColor: "rgba(54, 162, 235, 1)",
                    borderWidth: 1,
                },
            ],
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    display: true,
                },
            },
        },
    });
}

/**
 * Initialize the User Activity Chart (Line Chart)
 */
function initializeUserActivityChart() {
    const ctx = document.getElementById("userActivityChart").getContext("2d");

    // Simulated user activity data for demonstration (replace with API data if needed)
    const data = [5, 15, 10, 25];
    const labels = ["Jan", "Feb", "Mar", "Apr"];

    new Chart(ctx, {
        type: "line",
        data: {
            labels: labels,
            datasets: [
                {
                    label: "User Activity",
                    data: data,
                    borderColor: "rgba(255, 99, 132, 1)",
                    backgroundColor: "rgba(255, 99, 132, 0.2)",
                    borderWidth: 2,
                    fill: true,
                },
            ],
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    display: true,
                },
            },
        },
    });
}
