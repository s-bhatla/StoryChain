import { useState } from 'react'
import { Link, useNavigate } from "react-router-dom";
import { useParams } from "react-router-dom";


import '../App.css'

function Results() {
    const { roomID } = useParams();
    const [stories, setStories] = useState([]);
    const [currentStory, setCurrentStory] = useState(0);
    const [loading, setLoading] = useState(true);
    const navigate = useNavigate();

    useEffect(() => {
        const fetchStories = async () => {
            try {
                const response = await fetch(`/get-final-stories/${roomID}`, {
                    method: 'GET',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });
                const data = await response.json();
                console.log("Stories MAP: ", data)
                setStories(data);
                setLoading(false);
            } catch (error) {
                console.error(error);
            }
        };
        fetchStories();
    }, [roomID]);

    const handleNextStory = () => {
        setCurrentStory(currentStory + 1);
    };

    if (loading) {
        return <div>Loading...</div>;
    }

    if (Object.keys(stories).length === 0) {
        return <div>No stories found.</div>;
    }

    return (
        <div className="min-h-screen w-full bg-gray-50 text-gray-900 flex flex-col items-center justify-center">
            <h1>Stories</h1>
            {Object.keys(stories).map((username, index) => (
                <div key={index}>
                    <h2>{username}</h2>
                    <p>
                        {stories[username].map((storyLine, index) => (
                            <span key={index}>{storyLine} </span>
                        ))}
                    </p>
                </div>
            ))}
            {currentStory < stories.length - 1 && (
                <button onClick={handleNextStory}>Next Story</button>
            )}
        </div>
    )
}

export default Results
