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
    let config={};

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

    function loadConfig(){
        return $.get('/api/config',function(data){
            config=data;
        });
    }

    function showVoting() {
        breedSearch.style.display = 'none';
        breedInfo.style.display = 'none';
        votingButtons.style.display = 'flex';
        catContainer.style.display = 'grid';
        catContainer.innerHTML = '';
        getRandomCats(1);
    }

    function showBreeds() {
        breedSearch.style.display = 'block';
        breedInfo.style.display = 'block';
        votingButtons.style.display = 'none';
        catContainer.style.display = 'grid';
        catContainer.innerHTML = '';
        populateBreeds();
    }

    function showFavs() {
        breedSearch.style.display = 'none';
        breedInfo.style.display = 'none';
        votingButtons.style.display = 'none';
        catContainer.style.display = 'grid';
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
            const response = await fetch(`/api/random-cat?breed_ids=${breed.id}&count=5`);
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
        }, 3000);
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

    async function toggleFavorite(cat) {
        try {
            const isFav = await isFavorite(cat.id);
            let response;
    
            if (isFav) {
                response = await fetch(`/api/favorites/${cat.id}`, {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json',
                        'x-api-key': config.catapi_key
                    },
                });
            } else {
                response = await fetch('/api/favorites', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'x-api-key': config.catapi_key
                    },
                    body: JSON.stringify({
                        image_id: cat.id,
                        sub_id: userId
                    }),
                });
            }
    
            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || 'Failed to update favorite');
            }
    
            const result = await response.json();
            console.log(isFav ? 'Favorite removed:' : 'Favorite added:', result);
            updateHeartButton(cat.id, !isFav);
        } catch (error) {
            console.error('Error toggling favorite:', error);
            alert(error.message || 'An error occurred. Please try again later.');
        }
    }
    
    async function isFavorite(catId) {
        try {
            const response = await fetch(`/api/favorites?sub_id=${userId}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            
            if (data.error) {
                console.error('API error:', data.error);
                return false;
            }
            
            if (!Array.isArray(data)) {
                console.error('Unexpected API response:', data);
                return false;
            }
            
            return data.some(fav => fav.image_id === catId);
        } catch (error) {
            console.error('Error checking favorite status:', error);
            return false;
        }
    }
    
    function updateHeartButton(favoriteId, isFav, imageId) {
        const heartButton = document.querySelector(`.cat-item[data-id="${imageId}"] .heart-button`);
        if (heartButton) {
            heartButton.innerHTML = isFav ? '❤️' : '♡';
            heartButton.classList.toggle('favorited', isFav);
            heartButton.onclick = (event) => {
                event.stopPropagation();
                if (isFav) {
                    removeFavorite(favoriteId);
                } else {
                    toggleFavorite({ id: imageId });
                }
            };
        }
    }
    async function removeFavorite(favoriteId) {
        try {
            const response = await fetch(`/api/favorites/${favoriteId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                    'x-api-key': config.catapi_key
                },
            });
    
            const responseData = await response.json();
            console.log('Response:', responseData);
    
            if (!response.ok) {
                throw new Error(responseData.error || 'Failed to remove favorite');
            }
    
            console.log('Favorite removed:', favoriteId);
            // Remove the cat element from the DOM
            const catElement = document.querySelector(`.cat-item[data-id="${favoriteId}"]`);
            if (catElement) {
                catElement.remove();
            }
        } catch (error) {
            console.error('Error removing favorite:', error);
            alert(error.message || 'An error occurred. Please try again later.');
        }
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
            getRandomCats(1);
            event.stopPropagation();
            toggleFavorite(cat);
        });
    
        catElement.appendChild(img);
        catElement.appendChild(heartButton);
    
        // Check if it's a favorite and update the heart button accordingly
        isFavorite(cat.id).then(isFav => updateHeartButton(cat.id, isFav));
    
        return catElement;
    }
    
    async function displayFavoriteCats() {
        try {
            const response = await fetch(`/api/favorites?sub_id=${userId}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const data = await response.json();
            
            if (data.error) {
                throw new Error(data.error);
            }
            
            if (!Array.isArray(data)) {
                throw new Error('API response is not an array');
            }
    
            catContainer.innerHTML = '';
            for (const fav of data) {
                if (fav.image) {
                    const catElement = createCatElement(fav.image);
                    updateHeartButton(fav.id, true, fav.image.id);
                    catContainer.appendChild(catElement);
                } else {
                    console.warn('Favorite item does not contain image data:', fav);
                }
            }
        } catch (error) {
            console.error('Error fetching favorite cats:', error);
            alert('Failed to fetch favorite cats: ' + error.message);
        }
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
