function updateModels() {
    const brandSelect = document.getElementById('brand');
    const modelSelect = document.getElementById('model');
    const modelsData = JSON.parse(document.getElementById('models-data').textContent);

    brandSelect.addEventListener('change', function() {
        const brandId = this.value;
        modelSelect.innerHTML = '<option value="">Все модели</option>';
        
        if (!brandId) {
            modelSelect.disabled = true;
            return;
        }

        modelSelect.disabled = false;
        const filteredModels = modelsData.filter(model => model.brand_id == brandId);
        
        filteredModels.forEach(model => {
            const option = document.createElement('option');
            option.value = model.id;
            option.textContent = model.name;
            modelSelect.appendChild(option);
        });
    });
}

document.addEventListener('DOMContentLoaded', updateModels);