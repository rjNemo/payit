(() => {
  const form = document.querySelector("form");
  const button = document.querySelector("button");
  const qtyInput = document.querySelector("input#quantity");
  const message = document.querySelector("div#message");

  if (!form || !button || !qtyInput) {
    console.error("Missing required form elements");
    return;
  }

  const setMessage = (text, isError = true) => {
    if (!message) {
      return;
    }
    message.textContent = text;
    message.style.color = isError ? "#dc2626" : "#16a34a";
  };

  form.addEventListener("submit", async (event) => {
    event.preventDefault();

    const quantity = Number.parseInt(qtyInput.value, 10);
    if (!Number.isFinite(quantity) || quantity <= 0) {
      setMessage("Enter a quantity of at least 1.");
      qtyInput.focus();
      return;
    }

    try {
      button.disabled = true;
      setMessage("Contacting Stripe…", false);

      const response = await fetch("/api/checkout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ quantity }),
      });

      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || "Checkout request failed.");
      }

      const data = await response.json();
      if (!data || !data.url) {
        throw new Error("Checkout response missing redirect URL.");
      }

      setMessage("Redirecting to Stripe…", false);
      window.location.href = data.url;
    } catch (err) {
      console.error("Checkout failed", err);
      setMessage("Unable to start checkout. Please try again.");
      button.disabled = false;
    }
  });
})();
