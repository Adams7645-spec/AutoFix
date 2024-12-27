document.addEventListener('DOMContentLoaded', function() {
    loadCars();
});

function loadCars() {
    fetch('/api/cars')
        .then(response => response.json())
        .then(data => {
            data.forEach(car => {
                addCarToList(car);
            });
        })
        .catch(error => console.error('Ошибка:', error));
}

document.getElementById('carForm').addEventListener('submit', function(event) {
    event.preventDefault();

    const title = document.getElementById('title').value;
    const description = document.getElementById('description').value;
    const carNumber = document.getElementById('carNumber').value;
    const carBrand = document.getElementById('carBrand').value;
    const clientName = document.getElementById('clientName').value;
    const clientPhone = document.getElementById('clientPhone').value;

    const carRecord = {
        title: title,
        description: description,
        carNumber: carNumber,
        carBrand: carBrand,
        clientName: clientName,
        clientPhone: clientPhone
    };

    fetch('/api/cars', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(carRecord)
    })
    .then(response => response.json())
    .then(data => {
        addCarToList(data);
        document.getElementById('carForm').reset();
    })
    .catch(error => console.error('Ошибка:', error));
});

function addCarToList(car) {
    const tableBody = document.querySelector('#carList tbody');
    const row = document.createElement('tr');
    row.innerHTML = `
        <td class="editable" data-field="title">${car.title}</td>
        <td class="editable" data-field="description">${car.description}</td>
        <td class="editable" data-field="carNumber">${car.carNumber}</td>
        <td class="editable" data-field="carBrand">${car.carBrand}</td>
        <td class="editable" data-field="clientName">${car.clientName}</td>
        <td class="editable" data-field="clientPhone">${car.clientPhone}</td>
        <td>
            <div>
                <button onclick="deleteCar('${car.carNumber}')">Удалить</button>
                <button onclick="editCar('${car.carNumber}', this)">Редактировать</button>
            </div>
        </td>
    `;
    tableBody.appendChild(row);
}

function editCar(carNumber, button) {
    const row = button.closest('tr');
    const cells = row.querySelectorAll('.editable');

    cells.forEach(cell => {
        const field = cell.dataset.field;
        const input = document.createElement(field === 'description' ? 'textarea' : 'input');
        input.type = 'text';
        input.value = cell.textContent;
        cell.textContent = '';
        cell.appendChild(input);
    });

    button.textContent = 'Сохранить';
    button.setAttribute('onclick', `saveCar('${carNumber}', this)`);
}

function saveCar(carNumber, button) {
    const row = button.closest('tr');
    const cells = row.querySelectorAll('.editable');
    const updatedCar = {};

    cells.forEach(cell => {
        const field = cell.dataset.field;
        updatedCar[field] = cell.querySelector('input, textarea').value;
        cell.textContent = updatedCar[field];
    });

    fetch(`/api/cars/${carNumber}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(updatedCar)
    })
    .then(response => {
        if (response.ok) {
            return response.json();
        } else {
            return response.json().then(err => {
                console.error('Ошибка:', err);
            });
        }
    })
    .then(data => {
        const tableBody = document.querySelector('#carList tbody');
        const rows = tableBody.getElementsByTagName('tr');
        for (let i = 0; i < rows.length; i++) {
            if (rows[i].cells[2].innerText === carNumber) {
                rows[i].cells[0].innerText = data.title;
                rows[i].cells[1].innerText = data.description;
                rows[i].cells[3].innerText = data.carBrand;
                rows[i].cells[4].innerText = data.clientName;
                rows[i].cells[5].innerText = data.clientPhone;
                break;
            }
        }
    })
    .catch(error => console.error('Ошибка:', error));
    
    button.textContent = 'Редактировать';
    button.setAttribute('onclick', `editCar('${carNumber}', this)`);
}

function deleteCar(carNumber) {
    fetch(`/api/cars/${carNumber}`, {
        method: 'DELETE'
    })
    .then(response => {
        if (response.ok) {
            const tableBody = document.querySelector('#carList tbody');
            const rows = tableBody.getElementsByTagName('tr');
            for (let i = 0; i < rows.length; i++) {
                if (rows[i].cells[2].innerText === carNumber) {
                    rows[i].remove();
                    break; 
                }
            }
        } else {
            return response.json().then(err => {
                console.error('Ошибка:', err);
            });
        }
    })
    .catch(error => console.error('Ошибка:', error));
}
