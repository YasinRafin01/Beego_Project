document.addEventListener('DOMContentLoaded', () => {
    const votingLink = document.getElementById('voting');
    const breedsLink = document.getElementById('breeds');
    const favsLink = document.getElementById('favs');
    const breedSearch = document.getElementById('breed-search');
    const breedInput = document.getElementById('breed-input');
    const breedClose = document.getElementById('breed-close');
    const breedList = document.getElementById('breed-list');
    const catContainer = document.getElementById('cat-container');
    const breedInfo = document.getElementById('breed-info');
    const votingButtons = document.querySelector('.voting-buttons');
    const gridViewButton = document.querySelector('.grid-view');
    const listViewButton = document.querySelector('.list-view');

    let allBreeds = [];

    // Generate a unique user ID
    const userId = localStorage.getItem('userId') || `user-${Date.now()}`;
    localStorage.setItem('userId', userId);

    // Event listeners
    votingLink.addEventListener('click', showVoting);
    breedsLink.addEventListener('click', showBreeds);
    favsLink.addEventListener('click', showFavs);
    breedClose.addEventListener('click', toggleBreedSearch);
    gridViewButton.addEventListener('click', () => setView('grid'));
    listViewButton.addEventListener('click', () => setView('list'));
    breedInput.addEventListener('input', handleBreedSearch);
    breedList.addEventListener('click', handleBreedSelection);

    function showVoting() {
        breedSearch.style.display = 'none';
        breedInfo.style.display = 'none';
        votingButtons.style.display = 'flex';
        catContainer.style.display = 'grid'; // Display cat container
        // Clear any existing cat content
        catContainer.innerHTML = '';
        getRandomCats(1);
    }
    
    function showBreeds() {
        breedSearch.style.display = 'block';
        breedInfo.style.display = 'block';
        votingButtons.style.display = 'none';
        catContainer.style.display = 'grid'; // Display cat container
        // Clear any existing cat content
        catContainer.innerHTML = '';
        populateBreeds();
    }
    
    function showFavs() {
        breedSearch.style.display = 'none';
        breedInfo.style.display = 'none';
        votingButtons.style.display = 'none';
        catContainer.style.display = 'grid'; // Display cat container
        // Clear any existing cat content
        catContainer.innerHTML = '';
        displayFavoriteCats();
    }

    function toggleBreedSearch() {
        breedList.classList.toggle('show');
        breedInput.value = '';
        displayBreeds(allBreeds);
    }

    function setView(view) {
        if (view === 'grid') {
            catContainer.className = 'grid-view';
            gridViewButton.classList.add('active');
            listViewButton.classList.remove('active');
        } else {
            catContainer.className = 'list-view';
            listViewButton.classList.add('active');
            gridViewButton.classList.remove('active');
        }
    }

    function getRandomCats(count, breedId = '') {
        let url = `/api/random-cat?count=${count}`;
        if (breedId) {
            url += `&breed_ids=${breedId}`;
        }
        fetch(url)
            .then(response => response.json())
            .then(cats => {
                displayCats(cats);
                if (cats.length > 0) {
                    setupVotingButtons(cats[0].id);
                }
            })
            .catch(error => console.error('Error fetching random cats:', error));
    }

    async function populateBreeds() {
        try {
            const response = await fetch('/api/breeds');
            allBreeds = await response.json();
            displayBreeds(allBreeds);
            if (allBreeds.length > 0) {
                showBreedInfo(allBreeds[0]);
            }
        } catch (error) {
            console.error('Error fetching breeds:', error);
        }
    }

    function displayCats(cats) {
        catContainer.innerHTML = '';
        cats.forEach(cat => {
            const catElement = createCatElement(cat);
            catContainer.appendChild(catElement);
        });
    }

    function createCatElement(cat) {
        const catElement = document.createElement('div');
        catElement.className = 'cat-item';
        catElement.dataset.id = cat.id;

        const img = document.createElement('img');
        img.src = cat.url;
        img.alt = 'Cat';

        const heartButton = document.createElement('button');
        heartButton.className = 'heart-button';
        heartButton.innerHTML = '❤️';
        heartButton.addEventListener('click', (event) => {
            event.stopPropagation();
            toggleFavorite(cat);
            //getRandomCats(1);
            showVoting()
        });

        catElement.appendChild(img);
        catElement.appendChild(heartButton);

        return catElement;
    }

    function displayBreeds(breeds) {
        breedList.innerHTML = '';
        breeds.forEach(breed => {
            const breedItem = document.createElement('div');
            breedItem.className = 'breed-item';
            breedItem.textContent = breed.name;
            breedItem.dataset.breedId = breed.id;
            breedList.appendChild(breedItem);
        });
    }

    async function showBreedInfo(breed) {
        breedInfo.style.display = 'block';
        document.getElementById('breed-name').textContent = `${breed.name} (${breed.origin}) ${breed.alt_names || ''}`;
        document.getElementById('breed-description').textContent = breed.description;
        document.getElementById('breed-wiki').href = breed.wikipedia_url;
    
        try {
            const response = await fetch(`/api/random-cat?breed_ids=${breed.id}&count=5`);  // Fetch 5 images
            const breedImages = await response.json();
            if (breedImages.length > 0) {
                const breedImagesContainer = document.getElementById('breed-images');
                breedImagesContainer.innerHTML = '';
    
                breedImages.forEach((img, index) => {
                    const imageElement = document.createElement('img');
                    imageElement.src = img.url;
                    imageElement.alt = breed.name;
                    imageElement.style.display = index === 0 ? 'block' : 'none';
                    breedImagesContainer.appendChild(imageElement);
                });
    
                createImageSlider(breedImagesContainer);
            }
        } catch (error) {
            console.error('Error fetching breed images:', error);
        }
    }
    
    function createImageSlider(container) {
        let currentIndex = 0;
        const images = container.querySelectorAll('img');
        
        setInterval(() => {
            images[currentIndex].style.display = 'none';
            currentIndex = (currentIndex + 1) % images.length;
            images[currentIndex].style.display = 'block';
        }, 3000);  // Change image every 3 seconds
    }

    function handleBreedSearch(e) {
        const searchTerm = e.target.value.toLowerCase();
        const filteredBreeds = allBreeds.filter(breed => 
            breed.name.toLowerCase().includes(searchTerm)
        );
        displayBreeds(filteredBreeds);
        breedList.classList.add('show');
    }

    function handleBreedSelection(e) {
        if (e.target.classList.contains('breed-item')) {
            const selectedBreedId = e.target.dataset.breedId;
            const selectedBreed = allBreeds.find(breed => breed.id === selectedBreedId);
            if (selectedBreed) {
                showBreedInfo(selectedBreed);
                breedInput.value = selectedBreed.name;
                breedList.classList.remove('show');
            }
        }
    }

    function toggleFavorite(cat) {
        let favorites = JSON.parse(localStorage.getItem('favoriteCats')) || [];
        const index = favorites.findIndex(favCat => favCat.id === cat.id);

        if (index === -1) {
            favorites.push(cat);
            localStorage.setItem('favoriteCats', JSON.stringify(favorites));
        } else {
            favorites.splice(index, 1);
            localStorage.setItem('favoriteCats', JSON.stringify(favorites));
        }

        updateHeartButton(cat.id);
    }

    function isFavorite(catId) {
        const favorites = JSON.parse(localStorage.getItem('favoriteCats')) || [];
        return favorites.some(cat => cat.id === catId);
    }

    function updateHeartButton(catId) {
        const heartButton = document.querySelector(`.cat-item[data-id="${catId}"] .heart-button`);
        if (heartButton) {
            heartButton.innerHTML = isFavorite(catId) ? '❤️' : '♡';
        }
    }

    function displayFavoriteCats() {
        const favorites = JSON.parse(localStorage.getItem('favoriteCats')) || [];
        catContainer.innerHTML = '';
        favorites.forEach(cat => {
            const catElement = createCatElement(cat);
            catContainer.appendChild(catElement);
        });
    }

    function setupVotingButtons(imageId) {
        const upVoteButton = document.createElement('button');
        upVoteButton.textContent = '👍';
        upVoteButton.addEventListener('click', () => vote(imageId, 1));

        const downVoteButton = document.createElement('button');
        downVoteButton.textContent = '👎';
        downVoteButton.addEventListener('click', () => vote(imageId, -1));

        votingButtons.innerHTML = '';
        votingButtons.appendChild(upVoteButton);
        votingButtons.appendChild(downVoteButton);
    }

    async function vote(imageId, value) {
        try {
            const response = await fetch('/api/votes', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    image_id: imageId,
                    sub_id: userId,
                    value: value
                }),
            });

            if (response.ok) {
                console.log(`Vote recorded: ${value}`);
                getRandomCats(1);
            } else {
                console.error('Failed to record vote');
            }
        } catch (error) {
            console.error('Error recording vote:', error);
        }
    }

    async function getVotingHistory() {
        try {
            const response = await fetch(`/api/votes?sub_id=${userId}`);
            const votingHistory = await response.json();
            console.log('Voting History:', votingHistory);
        } catch (error) {
            console.error('Error fetching voting history:', error);
        }
    }

    // Initial load
    showVoting();
    getVotingHistory();
});
