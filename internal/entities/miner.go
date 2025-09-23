package entities

import (
    "image/color"
    "math"
    "math/rand"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/vector"
)

type Miner struct {
    BaseEntity
    MiningTimer    int
    MiningInterval int 
    Output         chan *Number
    MiningRange    int 
}

func NewMiner(x, y float64) *Miner {
    return &Miner{
        BaseEntity: BaseEntity{
            X:      x,
            Y:      y,
            Width:  32,
            Height: 32,
        },
        MiningTimer:    0,
        MiningInterval: 180, // 3 seconds at 60 FPS
        Output:         make(chan *Number, 10),
        MiningRange:    determineMiningRange(x, y),
    }
}

func (m *Miner) Update() {
    m.MiningTimer++
    
    if m.MiningTimer >= m.MiningInterval {
        m.MiningTimer = 0
        m.mine()
    }
}

func (m *Miner) Draw(screen *ebiten.Image, offsetX, offsetY float64, zoom float64) {
    screenX := float32((m.X + offsetX) * zoom)
    screenY := float32((m.Y + offsetY) * zoom)
    size := float32(m.Width * zoom)
    
    if size < 4 {
        return
    }
    
    baseColor := color.RGBA{150, 100, 50, 255}
    vector.DrawFilledRect(screen, screenX, screenY, size, size, baseColor, false)
    
    drillColor := color.RGBA{100, 100, 100, 255}
    drillSize := size * 0.6
    drillX := screenX + (size-drillSize)/2
    drillY := screenY + (size-drillSize)/2
    vector.DrawFilledCircle(screen, drillX+drillSize/2, drillY+drillSize/2, 
        drillSize/2, drillColor, false)
    
    progress := float32(m.MiningTimer) / float32(m.MiningInterval)
    if progress > 0.8 { 
        animColor := color.RGBA{255, 255, 100, 200}
        animRadius := size * 0.3 * (1 + (progress-0.8)*5)
        vector.StrokeCircle(screen, screenX+size/2, screenY+size/2, 
            animRadius, 2, animColor, false)
    }
    
    borderColor := color.RGBA{200, 150, 100, 255}
    vector.StrokeRect(screen, screenX, screenY, size, size, 2, borderColor, false)
}

func (m *Miner) mine() {
    distanceFromOrigin := math.Sqrt(m.X*m.X + m.Y*m.Y) / 100 
    
    var value int
    if distanceFromOrigin < 2 {
        value = rand.Intn(10) + 1
    } else if distanceFromOrigin < 5 {
        value = rand.Intn(100) + 1
    } else {
        if rand.Float64() < 0.7 {
            value = generatePrimeInRange(10, int(distanceFromOrigin*50))
        } else {
            value = rand.Intn(int(distanceFromOrigin*20)) + 10
        }
    }
    
    // Create number at miner location with slight offset
    offsetX := (rand.Float64() - 0.5) * 20
    offsetY := (rand.Float64() - 0.5) * 20
    
    number := NewNumber(m.X+m.Width/2+offsetX, m.Y+m.Height/2+offsetY, value)
    
    // Try to send to output channel (non-blocking)
    select {
    case m.Output <- number:
        // Number sent successfully
    default:
        // Channel full, drop the number
    }
}

func (m *Miner) GetMinedNumber() *Number {
    select {
    case number := <-m.Output:
        return number
    default:
        return nil
    }
}

func determineMiningRange(x, y float64) int {
    distance := math.Sqrt(x*x + y*y)
    return int(distance/100) + 1
}

func generatePrimeInRange(min, max int) int {
    if min > max || min < 2 {
        return 2
    }
    
    // Simple approach: generate random numbers and check if prime
    for attempts := 0; attempts < 100; attempts++ {
        candidate := rand.Intn(max-min+1) + min
        if isPrime(candidate) {
            return candidate
        }
    }
    
    // Fallback: return next prime after min
    for candidate := min; candidate <= max; candidate++ {
        if isPrime(candidate) {
            return candidate
        }
    }
    
    return 2 // Ultimate fallback
}