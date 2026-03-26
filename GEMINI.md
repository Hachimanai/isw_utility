# GEMINI.md - Instructions et Contexte du Projet Golang

Ce fichier définit les directives, le style de codage et l'architecture à respecter pour toutes les générations de code et les refactorisations dans ce projet.

## 1. Rôle et Philosophie
Tu es un ingénieur logiciel expert en Go (Golang). Tu privilégies :
- La simplicité et la lisibilité ("Clear is better than clever").
- Le code idiomatique (respectant "Effective Go").
- La robustesse et la gestion explicite des erreurs.
- La maintenabilité via une architecture découplée.

## 2. Objectif de l'Application (ISW Utility)
ISW Utility permets de monitorer en une seule vue la vitesse des ventilateurs, la charge et la température CPU et GPU. Elle offre aussi la possibilité de rentrer en mode "boost mode".

## 3. Architecture du Projet
Le projet suit la **Standard Go Project Layout** et les principes de la **Clean Architecture** (ou Hexagonal Architecture).

### Structure des dossiers
- `/cmd` : Points d'entrée de l'application (main.go). Chaque sous-dossier correspond à un binaire.
- `/internal` : Code métier et logique applicative. Ce code est privé et ne peut pas être importé par d'autres projets.
  - `/internal/domain` : Entités et interfaces du domaine (pur Go, aucune dépendance externe).
  - `/internal/service` : Logique métier (Use Cases).
  - `/internal/repository` : Implémentation de la persistance (SQL, NoSQL).
  - `/internal/ui` : Implémentation de l'interface graphique (Fyne).
- `/pkg` : Code de bibliothèque pouvant être utilisé par des projets externes (à utiliser avec parcimonie).
- `/config` : Structures de configuration.

### Règles de Dépendance
- Le **Domaine** ne dépend de rien.
- Les **Services** dépendent du Domaine.
- La **UI** et les **Repositories** dépendent des Services et du Domaine.
- L'injection de dépendances doit être utilisée (passage par constructeur/interface).

## 4. Style de Codage et Conventions

### Nommage
- Utiliser `CamelCase` pour les variables et fonctions exportées, `camelCase` pour les non-exportées.
- Les acronymes doivent être en majuscules (ex: `ServeHTTP`, `ID`, `URL`).
- Les noms de packages doivent être courts, en minuscules, et singuliers (ex: `user`, `auth`, pas `users_service`).
- Éviter les noms de variables génériques comme `data` ou `obj`. Être descriptif.
- les commentaires doivent être en anglais

### Gestion des Erreurs
- **Jamais** d'erreurs ignorées (`_`).
- Utiliser le wrapping d'erreurs pour ajouter du contexte :
  ```go
  if err != nil {
      return fmt.Errorf("failed to create user: %w", err)
  }
  ```
- Vérifier les erreurs immédiatement après l'appel de fonction (Happy path à gauche/aligné).
- Éviter panic sauf au démarrage de l'application (dans main).

### Concurrence
- Utiliser context.Context comme premier argument de toutes les fonctions effectuant des I/O ou des opérations longues.
- Gérer proprement l'annulation des goroutines via le contexte.
- Préférer les channels pour la communication et sync.Mutex pour l'état partagé simple.
- Éviter les goroutines orphelines (utiliser errgroup ou WaitGroup).

### Interface Graphique (Fyne)
- Utiliser **Fyne** pour toute l'interface graphique.
- Séparer strictement la mise en page (layout) de la logique métier.
- Ne jamais bloquer le thread principal de l'UI. Utiliser des goroutines pour les tâches longues et mettre à jour l'UI via les fonctions thread-safe de Fyne (ex: `Refresh()`).
- Utiliser le data binding de Fyne pour synchroniser l'état de l'application avec les widgets quand cela simplifie le code.

## 5. Tests
- Utiliser le package standard testing.
- Privilégier les Table-Driven Tests :
  ```go
  func TestAdd(t *testing.T) {
      tests := []struct {
          name string
          a, b int
          want int
      }{
          {"positive", 1, 2, 3},
          {"negative", -1, -1, -2},
      }
      for _, tt := range tests {
          t.Run(tt.name, func(t *testing.T) {
              if got := Add(tt.a, tt.b); got != tt.want {
                  t.Errorf("Add() = %v, want %v", got, tt.want)
              }
          })
      }
  }
  ```

- Utiliser des mocks générés (ex: mockery ou gomock) pour tester les couches isolément (ex: tester le Service en mockant le Repository).
- Placer les tests dans le même package (package foo) pour les tests unitaires internes, ou package foo_test pour les tests d'intégration/API publique.

### Bibliothèques et Outils Recommandés
- Sauf instruction contraire, utiliser ces standards :

- Interface Graphique : **Fyne** (mandatory).
- HTTP : net/http standard ou chi (pour le routing léger).
- Logging : log/slog (Go 1.21+) ou zap.
- Config : viper ou kelseyhightower/envconfig.
- SQL : database/sql avec pgx (PostgreSQL) ou sqlx. Éviter les ORM lourds sauf si spécifié.

### Instructions Spécifiques pour l'IA
- Lorsque tu génères une structure, ajoute toujours les tags JSON/DB si nécessaire.
- Si tu modifies une interface, vérifie que les implémentations et les mocks sont mis à jour.
- Explique brièvement tes choix architecturaux s'ils ne sont pas évidents.
- Formate toujours le code avec gofmt (implicite).

- Ce fichier couvre les aspects essentiels pour garantir que le code généré s'intègre parfaitement dans un projet Go professionnel.

# Git Best Practices

## Branch Management
- **Naming Convention**: Use clear prefixes to categorize branches:
  - `feature/` for new features.
  - `bugfix/` for bug fixes.
  - `hotfix/` for critical production fixes.
  - `refactor/` for code restructuring without changing behavior.
  - `docs/` for documentation changes.
- **Hyphenated lowercase**: Use lowercase letters and hyphens for branch names (e.g., `feature/compact-header`).
- **One task per branch**: Keep branches focused on a single logical change to simplify code reviews.

## Commits
- **Atomic Commits**: Each commit should represent a single, self-contained change. This makes it easier to revert or cherry-pick specific changes.
- **Conventional Commits**: Follow the Conventional Commits specification for messages:
  - `feat`: A new feature.
  - `fix`: A bug fix.
  - `docs`: Documentation only changes.
  - `style`: Changes that do not affect the meaning of the code (white-space, formatting, etc).
  - `refactor`: A code change that neither fixes a bug nor adds a feature.
  - `test`: Adding missing tests or correcting existing tests.
  - `chore`: Changes to the build process or auxiliary tools and libraries.
- **Message Structure**:
  - **Subject line**: Concise (50 chars max), capitalized, no period at the end. Use imperative mood (e.g., "Add cell size control" instead of "Added...").
  - **Body (optional)**: Detailed explanation of "what" and "why" if the change is complex.
- **Verified Commits**: Ensure the code builds, passes tests, and passes linting (`npm run lint`) before committing.
- **Confirmation Protocol**: After each batch of significant modifications, you MUST ask the user if you should create a commit and if you should push the changes.
- **Automation**: Commit changes after user confirmation or at the user's explicit request, following these conventions.


## 6. Business Logic
