function openModal(type) {
    document.getElementById(`modal-${type}`).style.display = "flex";
  }

  function closeModal(type) {
    document.getElementById(`modal-${type}`).style.display = "none";
  }

  function highlightAssociations(clickedItem) {
    // First, clear any existing highlights
    document.querySelectorAll(".highlight").forEach((el) => {
      el.classList.remove("highlight");
    });

    // Highlight the clicked item
    clickedItem.classList.add("highlight");

    // Gather associated IDs from data attributes
    const dataBuildings = clickedItem.dataset.buildings?.split(",") || [];
    console.log("dataBuildings", dataBuildings)
    const dataProductions =
      clickedItem.dataset.productions?.split(",") || [];
    const dataProducts = clickedItem.dataset.products?.split(",") || [];

    // Highlight associated items
    dataBuildings.forEach((id) => {
      const buildingEl = document.getElementById(id.trim());
      if (buildingEl) buildingEl.classList.add("highlight");
    });

    dataProductions.forEach((id) => {
      const productionEl = document.getElementById(id.trim());
      if (productionEl) productionEl.classList.add("highlight");
    });

    dataProducts.forEach((id) => {
      const productEl = document.getElementById(id.trim());
      if (productEl) productEl.classList.add("highlight");
    });
  }

  function createBuilding() {
    console.log("creating building");
    console.log(document.getElementById("building-productions").value);
    console.log(
      document.getElementById("building-productions").value.split(",")
    );
    console.log(
      document
        .getElementById("building-productions")
        .value.split(",")
        .map((id) => parseInt(id))
    );
    const buildingData = {
      name: document.getElementById("building-name").value,
      resourceCost: parseFloat(
        document.getElementById("building-cost").value
      ),
      buildTime: parseInt(
        document.getElementById("building-build-time").value
      ),
      productions: document
        .getElementById("building-productions")
        .value.split(",")
        .map((id) => parseInt(id)),
    };
    console.log(buildingData);
    fetch("/api/buildings", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(buildingData),
    })
      .then((response) => {
        console.log(response.status);
        if (response.status !== 201) {
          throw new Error(
            `Network response was not ok with status: ${response.status}`
          );
        }
      })
      .then(() => {
        closeModal("building");
        window.location.reload();
      })
      .catch((error) => {
        console.error("Error creating building:", error);
        alert("Failed to create building. Please try again.");
      });
  }

  function createProduction() {
    const inputs = document
      .getElementById("input-productions")
      .value.split(",")
      .map((item) => {
        const splitItem = item.split(":");
        const productId = parseInt(splitItem[0]);
        const amount = parseInt(splitItem[1]);
        return {
          productId,
          amount,
        };
      });
    const outputs = document
      .getElementById("output-productions")
      .value.split(",")
      .map((item) => {
        const splitItem = item.split(":");
        const productId = parseInt(splitItem[0]);
        const amount = parseInt(splitItem[1]);
        return {
          productId,
          amount,
        };
      });
    const productionData = {
      name: document.getElementById("production-name").value,
      cost: parseFloat(document.getElementById("production-cost").value),
      duration: parseInt(
        document.getElementById("production-duration").value
      ),
      inputs,
      outputs,
    };
    fetch("/api/production", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(productionData),
    })
      .then(() => {
        closeModal("production");
        window.location.reload();
      })
      .catch((error) => {
        console.error("Error creating production:", error);
        alert("Failed to create production. Please try again.");
      });
  }

  function createProduct() {
    const productData = {
      name: getElementById("product-name").value,
      value: getElementById("product-value").value,
    };

    fetch("/api/product", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(productData),
    })
      .then(() => {
        closeModal("product");
        window.location.reload();
      })
      .catch((error) => {
        console.error("Error creating production:", error);
        alert("Failed to create product. Please try again.");
      });
  }

  function saveFiles() {
    fetch("/config/storage", {
      method: "POST",
    })
      .then(() => {
        alert("succesfully saved files");
      })
      .catch((error) => {
        console.error("Error saving Files:", error);
        alert("Failed to save files");
      });
  }