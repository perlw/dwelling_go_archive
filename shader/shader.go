package shader

import (
	"dwelling/math/matrix"
	"dwelling/math/vector"
	"errors"
	"fmt"
	gl "github.com/chsc/gogl/gl33"
	"io/ioutil"
)

type ShaderProgram struct {
	name       string
	programId  gl.Uint
	uniformIds map[string]gl.Int
}

type AttribLocation struct {
	Position int
	Location string
}

func LoadShaderProgram(programName string, attribLocations []AttribLocation) (*ShaderProgram, error) {
	// Vertex shader
	vertexFile, err := ioutil.ReadFile(programName + ".vert")
	if err != nil {
		return nil, errors.New("shader: Could not find matching vertex shader for \"" + programName + "\"")
	}
	vertexSource := gl.GLString(string(vertexFile))
	defer gl.GLStringFree(vertexSource)

	vertexObj := gl.CreateShader(gl.VERTEX_SHADER)
	gl.ShaderSource(vertexObj, 1, &vertexSource, nil)
	gl.CompileShader(vertexObj)
	defer gl.DeleteShader(vertexObj)

	fmt.Println("Vertex log")
	printShaderLog(vertexObj)

	// Fragment shader
	fragmentFile, err := ioutil.ReadFile(programName + ".frag")
	if err != nil {
		return nil, errors.New("shader: Could not find matching fragment shader for \"" + programName + "\"")
	}
	fragmentSource := gl.GLString(string(fragmentFile))
	defer gl.GLStringFree(fragmentSource)

	fragmentObj := gl.CreateShader(gl.FRAGMENT_SHADER)
	gl.ShaderSource(fragmentObj, 1, &fragmentSource, nil)
	gl.CompileShader(fragmentObj)
	defer gl.DeleteShader(fragmentObj)

	fmt.Println("Fragment log")
	printShaderLog(fragmentObj)

	// Program
	program := gl.CreateProgram()

	// Bind attriblocations
	for _, attribLoc := range attribLocations {
		glLocString := gl.GLString(attribLoc.Location)
		gl.BindAttribLocation(program, gl.Uint(attribLoc.Position), glLocString)
		gl.GLStringFree(glLocString)
	}

	gl.AttachShader(program, vertexObj)
	gl.AttachShader(program, fragmentObj)

	gl.LinkProgram(program)

	fmt.Println("Program log")
	printProgramLog(program)

	return &ShaderProgram{
		name:       programName,
		programId:  program,
		uniformIds: make(map[string]gl.Int),
	}, nil
}

func printShaderLog(shader gl.Uint) {
	var length gl.Int
	gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)
	glString := gl.GLStringAlloc(gl.Sizei(length))
	defer gl.GLStringFree(glString)
	gl.GetShaderInfoLog(shader, gl.Sizei(length), nil, glString)
	fmt.Println(gl.GoString(glString))
}

func printProgramLog(shader gl.Uint) {
	var length gl.Int
	gl.GetProgramiv(shader, gl.INFO_LOG_LENGTH, &length)
	glString := gl.GLStringAlloc(gl.Sizei(length))
	defer gl.GLStringFree(glString)
	gl.GetProgramInfoLog(shader, gl.Sizei(length), nil, glString)
	fmt.Println(gl.GoString(glString))
}

func (program *ShaderProgram) Use() {
	gl.UseProgram(program.programId)
}

func (program *ShaderProgram) GetProgramId() gl.Uint {
	return program.programId
}

func (program *ShaderProgram) SetUniformInt(location string, value int) error {
	locId, err := program.getUniformLocation(location)
	if err != nil {
		return err
	}

	gl.Uniform1i(locId, gl.Int(value))
	return nil
}

func (program *ShaderProgram) SetUniformFloat(location string, value float32) error {
	locId, err := program.getUniformLocation(location)
	if err != nil {
		return err
	}

	gl.Uniform1f(locId, gl.Float(value))
	return nil
}

func (program *ShaderProgram) SetUniformVector2i(location string, value vector.Vector2i) error {
	locId, err := program.getUniformLocation(location)
	if err != nil {
		return err
	}

	glArray := value.ToGL()
	gl.Uniform2iv(locId, 1, &glArray[0])
	return nil
}

func (program *ShaderProgram) SetUniformVector3i(location string, value vector.Vector3i) error {
	locId, err := program.getUniformLocation(location)
	if err != nil {
		return err
	}

	glArray := value.ToGL()
	gl.Uniform3iv(locId, 1, &glArray[0])
	return nil
}

func (program *ShaderProgram) SetUniformVector2f(location string, value vector.Vector2f) error {
	locId, err := program.getUniformLocation(location)
	if err != nil {
		return err
	}

	glArray := value.ToGL()
	gl.Uniform2fv(locId, 1, &glArray[0])
	return nil
}

func (program *ShaderProgram) SetUniformVector3f(location string, value vector.Vector3f) error {
	locId, err := program.getUniformLocation(location)
	if err != nil {
		return err
	}

	glArray := value.ToGL()
	gl.Uniform3fv(locId, 1, &glArray[0])
	return nil
}

func (program *ShaderProgram) SetUniformVector4f(location string, value vector.Vector4f) error {
	locId, err := program.getUniformLocation(location)
	if err != nil {
		return err
	}

	glArray := value.ToGL()
	gl.Uniform4fv(locId, 1, &glArray[0])
	return nil
}

func (program *ShaderProgram) SetUniformMatrix(location string, value *matrix.Matrix) error {
	locId, err := program.getUniformLocation(location)
	if err != nil {
		return err
	}

	glArray := value.ToGL()
	gl.UniformMatrix4fv(locId, 1, gl.FALSE, &glArray[0])
	return nil
}

func (program *ShaderProgram) getUniformLocation(location string) (gl.Int, error) {
	if locId, ok := program.uniformIds[location]; ok {
		return locId, nil
	}

	glLocString := gl.GLString(location)
	locId := gl.GetUniformLocation(program.programId, glLocString)
	gl.GLStringFree(glLocString)

	if locId < 0 {
		return -1, errors.New("shader: Could not find location of \"" + location + "\" in \"" + program.name + "\"")
	}

	program.uniformIds[location] = locId

	return locId, nil
}
