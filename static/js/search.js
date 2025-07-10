document.addEventListener('DOMContentLoaded', function() {
    updateModels();
    
    const searchForm = document.getElementById('searchForm');
    const resultsDiv = document.getElementById('results');
    

    function validateForm() {
    let isValid = true;
    const brandId = document.getElementById('brand').value;
    const modelId = document.getElementById('model').value;
    const year = document.getElementById('year').value;
    const price = document.getElementById('price').value;
    
    
    document.querySelectorAll('.error-message').forEach(el => {
        el.style.display = 'none';
    });
    
    
    if (!brandId && !modelId && !year && !price) {
        document.getElementById('brand-error').style.display = 'block';
        isValid = false;
    }
    
    return isValid;
    }

    searchForm.addEventListener('submit', function(e) {
        e.preventDefault();

        if (!validateForm()) {
        return;
        }

        const brandId = document.getElementById('brand').value;
        const modelId = document.getElementById('model').value;
        const year = document.getElementById('year').value;
        const price = document.getElementById('price').value;
        
        if (!brandId && !modelId && !year && !price) {
            showError("Укажите хотя бы один критерий поиска");
            return;
        }

        if (year && isNaN(year)) {
            showError("Год должен быть числом");
            return;
        }

        if (price && isNaN(price)) {
            showError("Цена должна быть числом");
            return;
        }


        fetch(buildSearchUrl(brandId, modelId, year, price))
            .then(response => {
                if (!response.ok) throw new Error("Ошибка сервера");
                return response.json();
            })
            .then(cars => {
                if (!cars || cars.length === 0) {  
                    showMessage("Автомобили не найдены");
                } else {
                    renderResults(cars);
                }
            })
            .catch(error => {
                console.error("Ошибка:", error);
                showError("Произошла ошибка при поиске. Пожалуйста, попробуйте позже.");
            });
    });

    function buildSearchUrl(brandId, modelId, year, price) {
        const params = new URLSearchParams();
        if (brandId) params.append('brand', brandId);
        if (modelId) params.append('model', modelId);
        if (year) params.append('year', year);
        if (price) params.append('price', price);
        return `/search?${params.toString()}`;
    }

    function renderResults(cars) {
        resultsDiv.innerHTML = `
            <table>
                <thead>
                    <tr>
                        <th>Марка</th>
                        <th>Модель</th>
                        <th>Год</th>
                        <th>Цена</th>
                    </tr>
                </thead>
                <tbody>
                    ${cars.map(car => `
                        <tr>
                            <td>${car.brand_name}</td>
                            <td>${car.model_name}</td>
                            <td>${car.year}</td>
                            <td>${car.price.toLocaleString()} тыс. руб.</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;
    }

    function showMessage(text) {
        resultsDiv.innerHTML = `<div class="message">${text}</div>`;
    }

    function showError(text) {
        resultsDiv.innerHTML = `<div class="error">${text}</div>`;
    }
});