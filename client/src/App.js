import React, { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Register from './components/Register';
import Login from './components/Login';
import GameLobby from './components/GameLobby';
import Game from './components/Game';
import { initWebSocket } from './WebSocket';

function App() {
  useEffect(() => {
    initWebSocket();
  }, []);

  return (
    <Router>
      <div className="App">
        <Routes>
          <Route path="/register" element={<Register />} />
          <Route path="/login" element={<Login />} />
          <Route path="/lobby" element={<GameLobby />} />
          <Route path="/game/:id" element={<Game />} />
          <Route path="/" element={<Login />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;