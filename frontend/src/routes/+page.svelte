<script>
    import { onMount } from 'svelte';

    let images = [];
    let loading = true;
    let error = null;

    async function fetchImages() {
        try {
            loading = true;
            error = null;
            const response = await fetch('http://localhost:5174/api/images');
            
            if (!response.ok) {
                throw new Error('Failed to fetch images');
            }
            
            images = await response.json();
        } catch (err) {
            error = err.message;
            console.error('Error loading gallery:', err);
        } finally {
            loading = false;
        }
    }

    onMount(fetchImages);
</script>

<svelte:head>
    <title>AfroBase - Cultural Art Gallery</title>
</svelte:head>

<div class="container">
    <div class="header">
        <h1>🎨 AfroBase Gallery</h1>
        <p>Preserving African Heritage • Empowering Artists</p>
    </div>

    {#if loading}
        <div class="loading-state">
            <p>Loading artwork...</p>
        </div>
    {:else if error}
        <div class="error-state">
            <h3>⚠️ Unable to load gallery</h3>
            <p>{error}</p>
            <button on:click={fetchImages}>Try Again</button>
        </div>
    {:else if images.length === 0}
        <div class="empty-state">
            <h3>🎨 No artwork yet</h3>
            <p>Be the first to share your creativity with the community!</p>
        </div>
    {:else}
        <div class="gallery">
            {#each images as image}
                <div class="art-card">
                    <div class="art-image">
                        <img src={image.url} alt={image.title} />
                    </div>
                    <div class="art-info">
                        <div class="art-title">{image.title}</div>
                        <div class="art-description">{image.description}</div>
                        <div class="art-meta">
                            Shared on {new Date(image.upload_time * 1000).toLocaleDateString()}
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<style>
    * {
        margin: 0;
        padding: 0;
        box-sizing: border-box;
    }

    :global(body) {
        font-family: 'Ubuntu', Arial, sans-serif;
        line-height: 1.6;
        background: linear-gradient(135deg, #1a0e0a 0%, #2d1810 50%, #1a0e0a 100%);
        color: #f4f1eb;
        min-height: 100vh;
    }

    .container {
        max-width: 1200px;
        margin: 0 auto;
        padding: 20px;
    }

    .header {
        text-align: center;
        padding: 40px 0;
        background: radial-gradient(ellipse at center, rgba(218, 165, 32, 0.3) 0%, transparent 70%);
        border-radius: 20px;
        margin-bottom: 40px;
    }

    .header h1 {
        font-size: 3rem;
        font-weight: 700;
        margin-bottom: 10px;
        background: linear-gradient(45deg, #DAA520, #FF8C00, #CD853F);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
    }

    .header p {
        font-size: 1.2rem;
        color: #DAA520;
        margin-bottom: 20px;
    }

    .gallery {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
        gap: 30px;
        margin-top: 40px;
    }

    .art-card {
        background: rgba(255, 255, 255, 0.05);
        border-radius: 15px;
        border: 2px solid rgba(218, 165, 32, 0.2);
        overflow: hidden;
        transition: all 0.3s ease;
        backdrop-filter: blur(10px);
    }

    .art-card:hover {
        transform: translateY(-10px);
        border-color: #DAA520;
        box-shadow: 0 20px 40px rgba(218, 165, 32, 0.2);
    }

    .art-image {
        width: 100%;
        height: 250px;
        overflow: hidden;
    }

    .art-image img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.3s ease;
    }

    .art-card:hover .art-image img {
        transform: scale(1.05);
    }

    .art-info {
        padding: 25px;
    }

    .art-title {
        font-size: 1.4rem;
        font-weight: 600;
        color: #DAA520;
        margin-bottom: 10px;
    }

    .art-description {
        color: #f4f1eb;
        line-height: 1.6;
        margin-bottom: 15px;
    }

    .art-meta {
        color: rgba(255, 140, 0, 0.8);
        font-size: 0.9rem;
    }

    .empty-state, .loading-state, .error-state {
        text-align: center;
        padding: 80px 20px;
        color: rgba(218, 165, 32, 0.6);
    }

    .empty-state h3, .error-state h3 {
        font-size: 1.8rem;
        margin-bottom: 15px;
    }

    .error-state button {
        background: rgba(255, 140, 0, 0.2);
        color: #FF8C00;
        border: 2px solid rgba(255, 140, 0, 0.3);
        padding: 12px 25px;
        border-radius: 8px;
        cursor: pointer;
        margin-top: 20px;
        transition: all 0.3s ease;
    }

    .error-state button:hover {
        background: rgba(255, 140, 0, 0.3);
        border-color: #FF8C00;
    }

    @media (max-width: 768px) {
        .header h1 {
            font-size: 2rem;
        }
        
        .gallery {
            grid-template-columns: 1fr;
        }
    }
</style>