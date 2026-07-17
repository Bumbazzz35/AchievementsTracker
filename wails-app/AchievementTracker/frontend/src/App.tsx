import {useState, useEffect} from 'react';
import './App.css';
import {
    GetFullWorldState,
    StartWatching,
    StopWatching,
    SaveConfig,
    LoadConfig,
    SendTestEvent,
} from "../wailsjs/go/main/App";
import {EventsOn} from "../wailsjs/runtime/runtime";

function App() {
    const [worldName, setWorldName] = useState("");
    const [watching, setWatching] = useState(false);
    const [showConfig, setShowConfig] = useState(false);
    const [config, setConfig] = useState({minecraftPath: "", pollInterval: 5});
    const [events, setEvents] = useState<string[]>([]);

    useEffect(() => {
        EventsOn("advancement-event", (data: any) => {
            const eventText = data.type === "startup"
                ? `Запуск: мир "${data.worldName}" [${data.progressCurrent}/${data.progressTotal}]`
                : data.type === "new_advancement"
                    ? `Новое: ${data.title}`
                    : data.type === "criteria_update"
                        ? `Критерий: ${data.title} → ${(data.criteriaUpdates || []).join(", ")}`
                        : data.type === "progress_update"
                            ? `Обновление: ${data.title || ""}`
                            : `Событие: ${JSON.stringify(data)}`;
            setEvents(prev => [eventText || JSON.stringify(data), ...prev].slice(0, 50));
        });

        GetFullWorldState().then((state: any) => {
            if (state && state.WorldName) {
                setWorldName(state.WorldName);
            }
        });
        LoadConfig().then((cfg: any) => {
            if (cfg) {
                setConfig(cfg);
            }
        });
    }, []);

    const toggleWatching = () => {
        if (watching) {
            StopWatching();
            setWatching(false);
        } else {
            StartWatching().then(() => setWatching(true)).catch((err: any) => {
                alert("Ошибка: " + err);
            });
        }
    };

    const saveConfig = () => {
        SaveConfig(config).then(() => setShowConfig(false));
    };

    return (
        <div id="App">
            <h1>Minecraft Achievement Tracker</h1>
            {worldName && <p>Мир: {worldName}</p>}
            {!worldName && <p>Мир не найден. Укажи путь в настройках.</p>}

            <div>
                <button onClick={toggleWatching}>
                    {watching ? "Остановить" : "Отслеживать"}
                </button>
                <button onClick={() => setShowConfig(!showConfig)} style={{marginLeft: "10px"}}>
                    Настройки
                </button>
                <button onClick={() => SendTestEvent()} style={{marginLeft: "10px"}}>
                    Тест события
                </button>
            </div>

            {showConfig && (
                <div style={{marginTop: "20px"}}>
                    <p>Путь к .minecraft:</p>
                    <input
                        value={config.minecraftPath}
                        onChange={e => setConfig({...config, minecraftPath: e.target.value})}
                    />
                    <p>Интервал проверки (сек):</p>
                    <input
                        type="number"
                        value={config.pollInterval}
                        onChange={e => setConfig({...config, pollInterval: parseInt(e.target.value) || 5})}
                    />
                    <button onClick={saveConfig}>Сохранить</button>
                </div>
            )}

            <div style={{marginTop: "30px"}}>
                <h3>События:</h3>
                {events.length === 0 && <p>Нет событий</p>}
                <ul style={{textAlign: "left", maxHeight: "300px", overflowY: "auto"}}>
                    {events.map((e, i) => <li key={i}>{e}</li>)}
                </ul>
            </div>
        </div>
    )
}

export default App
