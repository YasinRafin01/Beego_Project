<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>The Cat API Clone</title>
    <link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
    <div class="container">
        <nav>
            <ul>
                <li><a href="#" id="voting"><span class="icon">↑↓</span> Voting</a></li>
                <li><a href="#" id="breeds"><span class="icon">🔍</span> Breeds</a></li>
                <li><a href="#" id="favs"><span class="icon">❤️</span> Favs</a></li>
            </ul>
            <div class="view-toggle">
                <button class="grid-view active">☰</button>
                <button class="list-view">≡</button>
            </div>
        </nav>

        <div id="breed-search" class="breed-search-container">
            <div class="breed-search-input">
                <input type="text" id="breed-input" placeholder="Enter a breed">
                <button id="breed-close">Select</button>
            </div>
            <div id="breed-list" class="breed-list"></div>
        </div>

        <div id="cat-container" class="grid-view">
            <!-- Cat images will be dynamically inserted here -->
        </div>

        <div id="breed-info">
            <div id="breed-images"></div>
            <h2 id="breed-name"></h2>
            <p id="breed-description"></p>
            <p>Source: <a href="#" id="breed-wiki" target="_blank">WIKIPEDIA</a></p>
        </div>

        <div class="voting-buttons" style="display: none;">
            <button class="like">👍</button>
            <button class="dislike">👎</button>
        </div>
    </div>

    <script src="/static/js/main.js"></script>
</body>
</html>
