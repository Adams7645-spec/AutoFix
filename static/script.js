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

    // Отправка данных на сервер
    fetch('/api/cars', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(carRecord)
    })
    .then(response => response.json())
    .then(data => {
        // Обновление списка автомобилей
        addCarToList(data);
        // Очистка формы
        document.getElementById('carForm').reset();
    })
    .catch(error => console.error('Ошибка:', error));
});

function addCarToList(car) {
    const tableBody = document.querySelector('#carList tbody');
    const row = document.createElement('tr');
    row.innerHTML = `
        <td>${car.title}</td>
        <td>${car.description}</td>
        <td>${car.carNumber}</td>
        <td>${car.carBrand}</td>
        <td>${car.clientName}</td>
        <td>${car.clientPhone}</td>
        <td>
            <div>
                <button onclick="deleteCar('${car.carNumber}')">Удалить</button>
                <button onclick="editCar('${car.carNumber}')">Редактировать</button>
            </div>
        </td>
    `;
    tableBody.appendChild(row);
}

function editCar(carNumber) {
    fetch('/api/cars/')
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
